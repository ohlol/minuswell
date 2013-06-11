package main

import (
	"fmt"
)

type PipeOutput struct {
	Formatter StringFormatter
}

func (p *PipeOutput) Emit(event Event) {
	fmt.Println(string(p.Formatter.Format(event)))
}
