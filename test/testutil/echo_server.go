package testutil

import (
	"bufio"
	"fmt"
	"gopolar/internal/core"
	"net"
	"os"
)

type EchoServer struct {
	name     string
	listener net.Listener
	prefix   string
}

// setup a echo server that replies same message with a prefix
func NewEchoServer(port uint64, prefix string) *EchoServer {
	p := ":" + fmt.Sprint(port)
	listener, err := net.Listen("tcp", p)
	ret := &EchoServer{
		name:     "[server" + p + "] ",
		listener: listener,
		prefix:   prefix,
	}
	if err != nil {
		core.Debugln(ret.name+"failed to create listener, err:", err)
		os.Exit(1)
	}
	core.Debugf(ret.name+"listening on %s, prefix: %s\n", listener.Addr(), prefix)
	go ret.run()
	return ret
}

func (es *EchoServer) run() {
	for {
		conn, err := es.listener.Accept()
		if err != nil {
			break
		}
		core.Debugln(es.name + "connected to " + conn.RemoteAddr().String())
		go es.handleConnection(conn)
	}
	core.Debugln(es.name + "quit")
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
			core.Debugln(es.name + "disconnected " + conn.RemoteAddr().String())
			break
		}
		core.Debugf(es.name+"request: %s", bytes)
		line := fmt.Sprintf("%s%s", es.prefix, bytes)
		core.Debugf(es.name+"response: %s", line)
		conn.Write([]byte(line))
	}
}
