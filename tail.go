package main

import (
	"github.com/ActiveState/tail"
	"github.com/howeyc/fsnotify"
	"log"
	"path/filepath"
)

type TailedFileLine struct {
	Filename string
	Type     string
	Tags     []string
	Fields   map[string]interface{}
	Line     *tail.Line
}

type TailedFile struct {
	Path    string
	Type    string
	Tags    []string
	Fields  map[string]interface{}
	Channel chan *TailedFileLine
}

func (t *TailedFile) Watch() {
	tl, _ := tail.TailFile(t.Path, tail.Config{Follow: true, ReOpen: true})
	log.Printf("Tailing file: %s\n", tl.Filename)

	for line := range tl.Lines {
		t.Channel <- &TailedFileLine{
			Filename: tl.Filename,
			Type:     t.Type,
			Tags:     t.Tags,
			Fields:   t.Fields,
			Line:     line,
		}
	}
}

func SetupWatcher(path string, config FilesConfig, ch chan *TailedFileLine) {
	tf := TailedFile{
		Path:    path,
		Type:    config.Type,
		Tags:    config.Tags,
		Fields:  config.Fields,
		Channel: ch,
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
