package main

import (
	"time"
)

type Event struct {
	Source     string                 `json:"@source"`
	Type       string                 `json:"@type"`
	Tags       []string               `json:"@tags,omitempty"`
	Fields     map[string]interface{} `json:"@fields,omitempty"`
	Timestamp  time.Time              `json:"@timestamp"`
	SourceHost string                 `json:"@source_host"`
	SourcePath string                 `json:"@source_path"`
	Message    string                 `json:"@message"`
	Formatter  FormatFunc             `json:"-"`
}
