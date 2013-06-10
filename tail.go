package main

import (
	"github.com/ActiveState/tail"
	"github.com/howeyc/fsnotify"
	"log"
	"path/filepath"
)

type TailedFileLine struct {
	Filename string
	Line     *tail.Line
}

type TailedFile struct {
	Path    string
	Channel chan *TailedFileLine
}

func (t *TailedFile) Watch() {
	tl, _ := tail.TailFile(t.Path, tail.Config{Follow: true, ReOpen: true})
	log.Printf("Tailing file: %s\n", tl.Filename)
	for line := range tl.Lines {
		t.Channel <- &TailedFileLine{Filename: tl.Filename, Line: line}
	}
}

func SetupWatcher(path string, ch chan *TailedFileLine) {
	tf := TailedFile{Path: path, Channel: ch}
	tf.Watch()
}

func WatchDir(path string, ch chan *TailedFileLine) {
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
				SetupWatcher(ev.Name, ch)
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
