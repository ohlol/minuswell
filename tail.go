package main

import (
	"github.com/ActiveState/tail"
	"github.com/howeyc/fsnotify"
	"log"
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
}

func (t *TailedFile) Watch() {
	tl, _ := tail.TailFile(t.Path, tail.Config{Follow: true, ReOpen: true})
	t.Logger.Printf("Tailing file: %s\n", tl.Filename)

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
		default:
			log.Printf("Buffer full while sending for: %s\n", tl.Filename)
		}
	}
}

func SetupWatcher(path string, config FilesConfig, logger *log.Logger, ch chan *TailedFileLine) {
	var tf TailedFile

	tf = TailedFile{
		Path:    path,
		Type:    config.Type,
		Tags:    config.Tags,
		Fields:  config.Fields,
		Channel: ch,
		Logger:  logger,
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

func WatchDirMask(path string, config FilesConfig, logger *log.Logger, ch chan *TailedFileLine) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		watcher.Close()
		return
	}

	err = watcher.Watch(filepath.Dir(path))
	if err != nil {
		log.Fatal(err)
		watcher.Close()
		return
	}

	for {
		select {
		case ev := <-watcher.Event:
			matched, err := regexp.MatchString(path, ev.Name)
			if err != nil {
				log.Println("error:", err)
			} else if matched {
				if ev.IsCreate() {
					SetupWatcher(ev.Name, config, logger, ch)
				} else if ev.IsDelete() {
					logger.Printf("file deleted: %s\n", ev.Name)
				} else if ev.IsRename() {
					logger.Printf("file renamed: %s\n", ev.Name)
				}
			}
		case err := <-watcher.Error:
			log.Println("error:", err)
			watcher.Close()
			return
		}
	}

	watcher.Close()
	return
}
