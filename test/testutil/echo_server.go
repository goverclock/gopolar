package testutil

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/goverclock/gopolar/internal/core"
)

type EchoServer struct {
	Name    string
	Prefix  string
	TotEcho uint64

	mu       sync.Mutex
	listener net.Listener
}

// setup a echo server that replies same message with a prefix
func NewEchoServer(port uint64, prefix string) *EchoServer {
	p := ":" + fmt.Sprint(port)
	listener, err := net.Listen("tcp", p)
	ret := &EchoServer{
		Name:     "[server" + p + "] ",
		listener: listener,
		Prefix:   prefix,
	}
	if err != nil {
		core.Debugln(ret.Name+"failed to create listener, err:", err)
		os.Exit(1)
	}
	core.Debugf(ret.Name+"listening on %s, prefix: %s\n", listener.Addr(), prefix)
	go ret.run()
	return ret
}

func (es *EchoServer) run() {
	for {
		conn, err := es.listener.Accept()
		if err != nil {
			break
		}
		core.Debugln(es.Name + "connected to " + conn.RemoteAddr().String())
		go es.handleConnection(conn)
	}
	core.Debugln(es.Name + "quit")
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
			core.Debugln(es.Name + "disconnected " + conn.RemoteAddr().String())
			break
		}
		core.Debugf(es.Name+"request: %s\\n", bytes[:len(bytes)-1])
		line := fmt.Sprintf("%s%s", es.Prefix, bytes)
		_, err = conn.Write([]byte(line))
		core.Debugf(es.Name+"response: %s\\n, err=%v", line[:len(line)-1], err)
		es.mu.Lock()
		es.TotEcho += uint64(len(line))
		es.mu.Unlock()
	}
}
