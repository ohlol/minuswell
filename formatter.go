package main

import (
	"encoding/json"
	"fmt"
)

type Formatter interface {
	Format(event Event) []byte
}

type FormatFunc func(event Event) []byte

type JsonFormatter struct{}
type RawFormatter struct{}
type StringFormatter struct{}

func (j *JsonFormatter) Format(event Event) []byte {
	ret, _ := json.Marshal(event)
	return ret
}

func (r *RawFormatter) Format(event Event) []byte {
	return []byte(event.Message)
}

func (s *StringFormatter) Format(event Event) []byte {
	return []byte(fmt.Sprintf("[%s] [%s] %s", event.Host, event.Timestamp, event.Message))
}
