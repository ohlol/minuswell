package main

type Formatter interface {
	JsonFormatter() []byte
	StringFormatter() string
}
