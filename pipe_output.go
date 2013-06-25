package main

import (
	"fmt"
	"log"
)

type PipeOutput struct {
	Host      string
	Formatter JsonFormatter
	Logger    *log.Logger
}

func (p *PipeOutput) Emit(event Event) {
	if event.Formatter != nil {
		fmt.Println(string(event.Formatter(event)))
	} else {
		fmt.Println(string(p.Formatter.Format(event)))
	}
}
