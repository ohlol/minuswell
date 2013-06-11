package main

import (
	"fmt"
	zmq "github.com/alecthomas/gozmq"
)

type ZmqOutput struct {
	Addresses []interface{}
	Formatter JsonFormatter
}

func (z *ZmqOutput) Emit(event Event) {
	ctx, _ := zmq.NewContext()
	socket, _ := ctx.NewSocket(zmq.PUSH)
	socket.SetSndTimeout(0)
	for _, addr := range z.Addresses {
		if err := socket.Connect(addr.(string)); err != nil {
			fmt.Println(err)
		}
	}

	socket.Send(z.Formatter.Format(event), 0)
}
