package main

import (
	"time"
)

type Event struct {
	Type      string
	Tags      []string
	Fields    map[string]interface{}
	Timestamp time.Time
	Host      string
	Path      string
	Message   string
	Formatter FormatFunc
}
