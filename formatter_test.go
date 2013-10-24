package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestJsonFormatter(t *testing.T) {
	flds := make(map[string]interface{})
	flds["fld1"] = 1

	e := Event{
		Type:      "test-event",
		Tags:      []string{"tag1", "tag2"},
		Fields:    flds,
		Timestamp: time.Now(),
		Host:      "localhost",
		Path:      "/path/here.log",
		Message:   "hello world",
	}
	enc, err := json.Marshal(e)
	if err != nil {
		t.Errorf("Problem marshalling event object to JSON: %s", err)
	}

	s := JsonFormatter{}
	if x := string(s.Format(e)); x != string(enc) {
		t.Errorf("JsonFormatter.Format(%v) = %v, want %v", e, string(x), string(enc))
	}
}

func TestStringFormatter(t *testing.T) {
	e := Event{
		Timestamp: time.Now(),
		Host:      "localhost",
		Message:   "hello world",
	}
	expectedString := fmt.Sprintf("[%s] [%s] %s", e.Host, e.Timestamp, e.Message)

	s := StringFormatter{}
	if x := string(s.Format(e)); x != expectedString {
		t.Errorf("StringFormatter.Format(%v) = %v, want %v", e, x, expectedString)
	}
}
