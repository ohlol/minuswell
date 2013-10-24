package main

import (
	"time"
)

type Event struct {
	Type       string                 `json:"type,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
	Timestamp  time.Time              `json:"@timestamp"`
	Host string                 `json:"host"`
	Path string                 `json:"path"`
	Message    string                 `json:"message"`
	Formatter  FormatFunc             `json:"-"`
}
