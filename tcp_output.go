package main

import (
	"fmt"
	"log"
	"net"
)

type TcpOutput struct {
	Host      string
	Port      int
	Formatter JsonFormatter
	connected bool
	conn      *net.Conn
}

func (t *TcpOutput) connect() error {
	log.Printf("[tcp] Connecting to: %s:%d", t.Host, t.Port)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.Host, t.Port))
	if err != nil {
		log.Printf("Could not connect to TCP: %s:%d", t.Host, t.Port)
		return err
	}

	t.conn = &conn
	t.connected = true
	return nil
}

func (t *TcpOutput) Emit(event Event) {
	if !t.connected {
		if err := t.connect(); err != nil {
			return
		}
	}

	fmt.Fprintf(*t.conn, string(t.Formatter.Format(event)))
}
