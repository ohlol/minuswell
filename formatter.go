package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

type Formatter interface {
	Format(event Event) []byte
}

type FormatFunc func(event Event) []byte

type JsonFormatter struct{}
type RawFormatter struct{}
type StringFormatter struct{}

func (j *JsonFormatter) Format(event Event) []byte {
	jsonEvent := make(map[string]interface{})
	jsonEvent["@timestamp"] = event.Timestamp
	jsonEvent["fqdn"] = event.Host
	jsonEvent["path"] = event.Path
	jsonEvent["message"] = strings.TrimRightFunc(event.Message, unicode.IsSpace)

	if event.Type != "" {
		jsonEvent["type"] = event.Type
	}
	if len(event.Tags) > 0 {
		jsonEvent["tags"] = event.Tags
	}

	for fn, fv := range event.Fields {
		jsonEvent[fn] = fv
	}

	ret, _ := json.Marshal(jsonEvent)
	return ret
}

func (r *RawFormatter) Format(event Event) []byte {
	return []byte(event.Message)
}

func (s *StringFormatter) Format(event Event) []byte {
	return []byte(fmt.Sprintf("[%s] [%s] %s", event.Host, event.Timestamp, strings.TrimRightFunc(event.Message, unicode.IsSpace)))
}
