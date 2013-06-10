package main

import (
	"time"
)

type Event struct {
	Source     string    `json:"@source"`
	Type       string    `json:"@type"`
	Tags       []string  `json:"@tags"`
	Fields     []string  `json:"@fields"`
	Timestamp  time.Time `json:"@timestamp"`
	SourceHost string    `json:"@source_host"`
	SourcePath string    `json:"@source_path"`
	Message    string    `json:"@message"`
}
