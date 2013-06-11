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
}

func (t *TcpOutput) Emit(event Event) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.Host, t.Port))
	if err != nil {
		log.Printf("Could not send to TCP: %s:%d", t.Host, t.Port)
		return
	}

	defer conn.Close()
	fmt.Fprintf(conn, string(t.Formatter.Format(event)))
}
