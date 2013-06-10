package main

import (
	"encoding/json"
	"fmt"
)

type Transport struct{}

func (t *Transport) JsonFormatter(event Event) string {
	j, err := json.Marshal(event)
	if err != nil {
		return ""
	}

	return string(j)
}

func (t *Transport) StringFormatter(event Event) string {
	return fmt.Sprintf("[%s] [%s] %s", event.SourceHost, event.Timestamp, event.Message)
}
