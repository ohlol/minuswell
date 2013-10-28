package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type tcpConnection struct {
	Host      string
	Port      int
	Conn      net.Conn
	connected bool
	timeout   time.Duration
}

type TcpOutput struct {
	Host      string
	Port      int
	Formatter JsonFormatter
	Logger    *log.Logger
	conn      tcpConnection
}

// defaultTimeout is the default connection timeout used by DialTimeout.
const (
	defaultTimeout  = 30
	initialWaitTime = 5
	maxWaitTime = 60
)

func (tc *tcpConnection) getConn() {
	var (
		err      error
		waitTime time.Duration
	)

	connectAddr := fmt.Sprintf("%s:%d", tc.Host, tc.Port)

	if tc.timeout == 0 {
		tc.timeout = defaultTimeout * time.Second
	}

	waitTime = initialWaitTime
	tc.Conn, err = net.DialTimeout("tcp", connectAddr, tc.timeout)
	for err != nil {
		log.Printf("[tcp] error connecting, retrying in %d seconds: %v\n", waitTime, err)
		time.Sleep(waitTime * time.Second)

		waitTime += 5
		if waitTime > maxWaitTime {
			waitTime = initialWaitTime
		}
		tc.Conn, err = net.DialTimeout("tcp", connectAddr, tc.timeout)
	}
}

func (t *TcpOutput) connect() {
	t.conn = tcpConnection{Host: t.Host, Port: t.Port}
	t.conn.getConn()
	t.conn.connected = true
}

func (t *TcpOutput) Emit(event Event) {
	var err error

	if !t.conn.connected {
		t.connect()
	}

	if event.Formatter != nil {
		_, err = fmt.Fprintf(t.conn.Conn, fmt.Sprintf("%s\n", string(event.Formatter(event))))
	} else {
		_, err = fmt.Fprintf(t.conn.Conn, fmt.Sprintf("%s\n", string(t.Formatter.Format(event))))
	}

	if err != nil {
		log.Printf("[tcp] Connect issue: %v\n", err)
		t.conn.connected = false
	}
}
