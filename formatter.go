package main

import (
	"encoding/json"
	"fmt"
)

type Formatter interface {
	Format(event Event) []byte
}

type JsonFormatter struct{}
type StringFormatter struct{}

func (j *JsonFormatter) Format(event Event) []byte {
	ret, _ := json.Marshal(event)
	return ret
}

func (s *StringFormatter) Format(event Event) []byte {
	return []byte(fmt.Sprintf("[%s] [%s] %s", event.SourceHost, event.Timestamp, event.Message))
}
