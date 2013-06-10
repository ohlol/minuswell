package main

import (
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var (
		opts struct {
			Config string `short:"c" description:"path to config file" required:"true"`
		}
	)

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	config, err := ReadConfig(opts.Config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	t := TcpTransport{Host: config.Transports["tcp"]["address"].(string), Port: int(config.Transports["tcp"]["port"].(float64))}
	fqdn, _ := os.Hostname()
	ch := make(chan *TailedFileLine)
	quit := make(chan bool)

	for fp, _ := range config.Files {
		go func(path string, c chan *TailedFileLine) {
			WatchDir(path, c)
		}(fp, ch)

		files, err := filepath.Glob(fp)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		for _, path := range files {
			go func(p string, c chan *TailedFileLine) {
				SetupWatcher(p, c)
			}(path, ch)
		}
	}

	go func() {
		for line := range ch {
			t.emit(Event{SourcePath: line.Filename, Timestamp: line.Line.Time, SourceHost: fqdn, Message: line.Line.Text})
		}
	}()

	<-quit
}
