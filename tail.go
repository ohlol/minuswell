package main

import (
	"github.com/ActiveState/tail"
	"github.com/howeyc/fsnotify"
	"log"
	"path/filepath"
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
}

func (t *TailedFile) Watch() {
	tl, _ := tail.TailFile(t.Path, tail.Config{Follow: true, ReOpen: true})
	log.Printf("Tailing file: %s\n", tl.Filename)

	for line := range tl.Lines {
		t.Channel <- &TailedFileLine{
			Filename:  tl.Filename,
			Type:      t.Type,
			Tags:      t.Tags,
			Fields:    t.Fields,
			Line:      line,
			Formatter: t.Formatter,
		}
	}
}

func SetupWatcher(path string, config FilesConfig, ch chan *TailedFileLine) {
	var tf TailedFile

	switch config.Format {
	case "json":
		formatter := JsonFormatter{}
		tf = TailedFile{
			Path:      path,
			Type:      config.Type,
			Tags:      config.Tags,
			Fields:    config.Fields,
			Channel:   ch,
			Formatter: formatter.Format,
		}
	case "string":
		formatter := StringFormatter{}
		tf = TailedFile{
			Path:      path,
			Type:      config.Type,
			Tags:      config.Tags,
			Fields:    config.Fields,
			Channel:   ch,
			Formatter: formatter.Format,
		}
	default:
		formatter := RawFormatter{}
		tf = TailedFile{
			Path:      path,
			Type:      config.Type,
			Tags:      config.Tags,
			Fields:    config.Fields,
			Channel:   ch,
			Formatter: formatter.Format,
		}
	}

	tf.Watch()
}

func WatchDir(path string, config FilesConfig, ch chan *TailedFileLine) {
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
			if ev.IsCreate() {
				SetupWatcher(ev.Name, config, ch)
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
