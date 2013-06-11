package main

import (
	zmq "github.com/alecthomas/gozmq"
	"log"
)

type ZmqOutput struct {
	Addresses []interface{}
	Formatter JsonFormatter
	connected bool
	socket    *zmq.Socket
}

func (z *ZmqOutput) connect() {
	ctx, _ := zmq.NewContext()
	socket, _ := ctx.NewSocket(zmq.PUSH)
	for _, addr := range z.Addresses {
		log.Printf("[zmq] Connecting to: %s\n", addr)
		socket.Connect(addr.(string))
	}

	z.socket = socket
	z.connected = true
}

func (z *ZmqOutput) Emit(event Event) {
	if !z.connected {
		z.connect()
	}

	z.socket.Send(z.Formatter.Format(event), 0)
}
