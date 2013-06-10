package main

import (
	"encoding/json"
	"fmt"
)

type Formatter interface {
	JsonFormatter() []byte
	StringFormatter() string
}

type Transport struct{}

func (t *Transport) JsonFormatter(event Event) []byte {
	j, err := json.Marshal(event)
	if err != nil {
		return nil
	}

	return j
}

func (t *Transport) StringFormatter(event Event) string {
	return fmt.Sprintf("[%s] [%s] %s", event.SourceHost, event.Timestamp, event.Message)
}
