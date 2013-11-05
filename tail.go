package main

import (
	"github.com/ActiveState/tail"
	"github.com/howeyc/fsnotify"
	sn "github.com/ohlol/shoenice"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type TailedFileLine struct {
	Filename  string
	Type      string
	Tags      []string
	Fields    map[string]interface{}
	Line      *tail.Line
	Formatter FormatFunc
}

type TailedFile struct {
	Path      string
	Type      string
	Tags      []string
	Fields    map[string]interface{}
	Channel   chan *TailedFileLine
	Formatter FormatFunc
	Logger    *log.Logger
	Stats     *sn.StatsInstance
}

func (t *TailedFile) Watch() {
	tl, _ := tail.TailFile(t.Path, tail.Config{Follow: true, Location: &tail.SeekInfo{0, os.SEEK_END}, Poll: true})
	t.Logger.Println("Tailing file:", tl.Filename)

	for line := range tl.Lines {
		select {
		case t.Channel <- &TailedFileLine{
			Filename:  tl.Filename,
			Type:      t.Type,
			Tags:      t.Tags,
			Fields:    t.Fields,
			Line:      line,
			Formatter: t.Formatter,
		}:
			t.Stats.Incr("tailed_lines")
		default:
			log.Println("Buffer full while sending for:", tl.Filename)
			t.Stats.Incr("buffer_full")
		}
	}
}

func SetupWatcher(path string, config FilesConfig, logger *log.Logger, ch chan *TailedFileLine, stats *sn.StatsInstance) {
	var tf TailedFile

	tf = TailedFile{
		Path:    path,
		Type:    config.Type,
		Tags:    config.Tags,
		Fields:  config.Fields,
		Channel: ch,
		Logger:  logger,
		Stats:   stats,
	}

	switch config.Format {
	case "json":
		formatter := JsonFormatter{}
		tf.Formatter = formatter.Format
	case "string":
		formatter := StringFormatter{}
		tf.Formatter = formatter.Format
	case "raw":
		formatter := RawFormatter{}
		tf.Formatter = formatter.Format
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("error: recovering file tail on %s: %s\n", path, r)
		}
	}()

	tf.Watch()
}

func WatchDirMask(path string, config FilesConfig, logger *log.Logger, ch chan *TailedFileLine, stats *sn.StatsInstance) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Printf("%s: %s\n", path, err)
		watcher.Close()
	}

	logger.Printf("Watching directory: %s\n", filepath.Dir(path))

	err = watcher.Watch(filepath.Dir(path))
	if err != nil {
		logger.Printf("%s: %s\n", path, err)
		watcher.Close()
	}

	for {
		select {
		case ev := <-watcher.Event:
			if ev != nil {
				matched, err := regexp.MatchString(path, ev.Name)
				if err != nil {
					log.Printf("%s: %s\n", path, err)
				} else if matched {
					if ev.IsCreate() {
						SetupWatcher(ev.Name, config, logger, ch, stats)
					}
				}
			}
		case err := <-watcher.Error:
			if err != nil {
				log.Println("error:", err)
				watcher.Close()
			}
		}
	}

	watcher.Close()
}
