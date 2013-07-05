package main

import (
	"fmt"
	"github.com/ohlol/go-flags"
	sn "github.com/ohlol/shoenice"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	GRAPHITE_PORT = 2003
	STATS_LISTEN_ADDR = "0.0.0.0"
	STATS_LISTEN_PORT = 5555
	STATS_PREFIX = "metrics"
	STATS_UPDATE_INTERVAL = 10
)

func formatFqdn() string {
        fqdn, _ := os.Hostname()
	splitName := strings.Split(fqdn, ".")

        for i, j := 0, len(splitName)-1; i < j; i, j = i+1, j-1 {
                splitName[i], splitName[j] = splitName[j], splitName[i]
        }

        return strings.Join(splitName, ".")
}

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

	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	config, err := ReadConfig(opts.Config)
	if err != nil {
		log.Fatal("Config parse error:", err)
	}

	if len(config.GraphiteHost) == 0 {
		log.Fatal("Did not specify Graphite host in config")
	}
	if config.GraphitePort == 0 {
		config.GraphitePort = 2003
	}
	if len(config.StatsListenAddr) == 0 {
		config.StatsListenAddr = STATS_LISTEN_ADDR
	}
	if config.StatsListenPort == 0 {
		config.StatsListenPort = STATS_LISTEN_PORT
	}
	if len(config.StatsPrefix) == 0 {
		config.StatsPrefix = strings.Join([]string{STATS_PREFIX, formatFqdn()}, ".")
	}
	if config.StatsUpdateInterval == 0 {
		config.StatsUpdateInterval = STATS_UPDATE_INTERVAL
	}

	stats := sn.NewStatsInstance()
	stats.RunServer(fmt.Sprintf("%s:%d", config.StatsListenAddr, config.StatsListenPort), config.StatsPrefix, config.StatsUpdateInterval, config.GraphiteHost, config.GraphitePort)

	if opts.ConfigDir != "" {
		cdir, err := os.Open(opts.ConfigDir)
		if err != nil {
			log.Println("Could not read config dir:", opts.ConfigDir)
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
			log.Fatal("Unknown output specified:", o)
		}

		stdoutLogger.Println("Setting up output:", o)

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
		go func(pth string, cfg FilesConfig) {
			WatchDirMask(pth, cfg, stdoutLogger, ch, stats)
		}(fp, config.Files[fp])

		files, err := filepath.Glob(fp)
		if err == nil {
			for _, path := range files {
				go func(pth string, cfg FilesConfig, c chan *TailedFileLine) {
					SetupWatcher(pth, cfg, stdoutLogger, c, stats)
				}(path, config.Files[fp], ch)
			}
		} else {
			log.Printf("%s: %s\n", fp, err)
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
