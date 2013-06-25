package main

import (
	"errors"
	"fmt"
	"github.com/ohlol/go-flags"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var (
		opts struct {
			Config    string   `short:"c" description:"Path to config file" required:"true" value-name:"FILE"`
			ConfigDir string   `short:"d" long:"config-dir" description:"Parse config files in dir" value-name:"DIR"`
			Output    []string `short:"o" description:"Which output to use (can specify multiple)" value-name:"OUTPUT"`
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

	if opts.ConfigDir != "" {
		cdir, err := os.Open(opts.ConfigDir)
		if err != nil {
			log.Printf("Could not read config dir: %s", opts.ConfigDir)
		} else {
			configs, _ := cdir.Readdirnames(-1)
			for _, cfg := range configs {
				pth := fmt.Sprintf("%s/%s", opts.ConfigDir, cfg)
				ncfg, err := ReadConfig(pth)
				if err != nil {
					log.Printf("Could not parse config file: %s (%s)\n", pth, err)
				} else {
					if ncfg.Outputs != nil {
						for k, v := range ncfg.Outputs {
							config.Outputs[k] = v
						}
					}
					if ncfg.Files != nil {
						for k, v := range ncfg.Files {
							config.Files[k] = v
						}
					}
				}
			}
		}
	}

	stdoutLogger := log.New(io.Writer(os.Stdout), "", log.LstdFlags)

	outputs := make([]Output, 0)
	for _, o := range opts.Output {
		if _, ok := outputsAvailable[o]; !ok {
			log.Fatal(errors.New(fmt.Sprintf("Unknown output specified.", o)))
			os.Exit(1)
		}

		stdoutLogger.Printf("Setting up %s output\n", o)

		switch o {
		case "pipe":
			outputs = append(outputs, &PipeOutput{Logger: stdoutLogger})
		case "tcp":
			outputs = append(outputs, &TcpOutput{Host: config.Outputs["tcp"]["address"].(string), Port: int(config.Outputs["tcp"]["port"].(float64)), Logger: stdoutLogger})
		case "zmq":
			outputs = append(outputs, &ZmqOutput{Addresses: config.Outputs["zmq"]["addresses"].([]interface{}), Logger: stdoutLogger})
		}
	}

	fqdn, _ := os.Hostname()
	ch := make(chan *TailedFileLine, 4096)
	quit := make(chan bool)

	for fp, _ := range config.Files {
		go func(pth string, cfg FilesConfig, c chan *TailedFileLine) {
			WatchDirMask(pth, cfg, stdoutLogger, c)
		}(fp, config.Files[fp], ch)

		files, err := filepath.Glob(fp)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		for _, path := range files {
			go func(pth string, cfg FilesConfig, c chan *TailedFileLine) {
				SetupWatcher(pth, cfg, stdoutLogger, c)
			}(path, config.Files[fp], ch)
		}
	}

	go func() {
		var event Event

		for line := range ch {
			event = Event{
				Source:     fmt.Sprintf("file://%s/%s", fqdn, line.Filename),
				Type:       line.Type,
				Tags:       line.Tags,
				Fields:     line.Fields,
				Timestamp:  line.Line.Time,
				SourceHost: fqdn,
				SourcePath: line.Filename,
				Message:    line.Line.Text,
			}
			if line.Formatter != nil {
				event.Formatter = line.Formatter
			}
			for _, t := range outputs {
				t.Emit(event)
			}
		}
	}()

	<-quit
}
