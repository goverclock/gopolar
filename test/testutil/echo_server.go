package testutil

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"testing"
)

type EchoServer struct {
	name     string
	listener net.Listener
	prefix   string
	t        *testing.T
}

// setup a echo server that replies same message with a prefix
func NewEchoServer(t *testing.T, port uint64, prefix string) *EchoServer {
	p := ":" + fmt.Sprint(port)
	listener, err := net.Listen("tcp", p)
	ret := &EchoServer{
		name:     "[server" + p + "] ",
		listener: listener,
		prefix:   prefix,
		t:        t,
	}
	if err != nil {
		t.Log(ret.name+"failed to create listener, err:", err)
		os.Exit(1)
	}
	t.Logf(ret.name+"listening on %s, prefix: %s\n", listener.Addr(), prefix)
	go ret.run()
	return ret
}

func (es *EchoServer) run() {
	for {
		conn, err := es.listener.Accept()
		if err != nil {
			break
		}
		es.t.Log(es.name + "connected to " + conn.RemoteAddr().String())
		go es.handleConnection(conn)
	}
	es.t.Log(es.name + "quit")
}

func (es *EchoServer) Quit() {
	es.listener.Close()
}

func (es *EchoServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			es.t.Log(es.name + "disconnected " + conn.RemoteAddr().String())
			break
		}
		es.t.Logf(es.name+"request: %s", bytes)
		line := fmt.Sprintf("%s%s", es.prefix, bytes)
		es.t.Logf(es.name+"response: %s", line)
		conn.Write([]byte(line))
	}
}
