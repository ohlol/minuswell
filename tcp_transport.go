package main

import (
	"fmt"
	"log"
	"net"
)

type TcpTransport struct {
	Host string
	Port int
	Transport
}

func (t *TcpTransport) emit(event Event) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.Host, t.Port))
	if err != nil {
		log.Printf("Could not send to TCP: %s:%d", t.Host, t.Port)
		return
	}

	defer conn.Close()
	fmt.Fprintf(conn, t.JsonFormatter(event))
}
