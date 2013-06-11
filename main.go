package main

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var (
		opts struct {
			Config string   `short:"c" description:"path to config file" required:"true"`
			Output []string `short:"o" description:"which output to use"`
		}
		outputsAvailable = map[string]bool{
			"pipe": true,
			"tcp":  true,
			"zmq":  true,
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

	outputs := make([]Output, 0)
	for _, o := range opts.Output {
		if _, ok := outputsAvailable[o]; !ok {
			log.Fatal(errors.New(fmt.Sprintf("Unknown output specified.", o)))
			os.Exit(1)
		}

		log.Printf("Setting up %s output\n", o)
		switch o {
		case "pipe":
			outputs = append(outputs, &PipeOutput{})
		case "tcp":
			outputs = append(outputs, &TcpOutput{Host: config.Outputs["tcp"]["address"].(string), Port: int(config.Outputs["tcp"]["port"].(float64))})
		case "zmq":
			outputs = append(outputs, &ZmqOutput{Addresses: config.Outputs["zmq"]["addresses"].([]interface{})})
		}
	}

	fqdn, _ := os.Hostname()
	ch := make(chan *TailedFileLine)
	quit := make(chan bool)

	for fp, _ := range config.Files {
		go func(pth string, cfg FilesConfig, c chan *TailedFileLine) {
			WatchDir(pth, cfg, c)
		}(fp, config.Files[fp], ch)

		files, err := filepath.Glob(fp)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		for _, path := range files {
			go func(pth string, cfg FilesConfig, c chan *TailedFileLine) {
				SetupWatcher(pth, cfg, c)
			}(path, config.Files[fp], ch)
		}
	}

	go func() {
		for line := range ch {
			for _, t := range outputs {
				t.Emit(Event{
					Source:     fmt.Sprintf("file://%s/%s", fqdn, line.Filename),
					Type:       line.Type,
					Tags:       line.Tags,
					Fields:     line.Fields,
					Timestamp:  line.Line.Time,
					SourceHost: fqdn,
					SourcePath: line.Filename,
					Message:    line.Line.Text,
				})
			}
		}
	}()

	<-quit
}
