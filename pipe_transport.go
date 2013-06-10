package main

import (
	"fmt"
)

type PipeTransport struct {
	Transport
}

func (pt *PipeTransport) emit(event Event) {
	fmt.Println(pt.StringFormatter(event))
}
