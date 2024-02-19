package testutil

import (
	"fmt"
	"net"
	"os"

	"github.com/goverclock/gopolar/internal/core"
)

type SilentServer struct {
	Name    string
	TotEcho uint64

	listener net.Listener
}

// setup a silent server that always disconnect immediately
func NewSilentServer(port uint64) *SilentServer {
	p := ":" + fmt.Sprint(port)
	listener, err := net.Listen("tcp", p)
	ret := &SilentServer{
		Name:     "[s-serv" + p + "] ",
		listener: listener,
	}
	if err != nil {
		core.Debugln(ret.Name+"failed to create listener, err:", err)
		os.Exit(1)
	}
	core.Debugf(ret.Name+"listening on %s\n", listener.Addr())
	go ret.run()
	return ret
}

func (ss *SilentServer) run() {
	for {
		conn, err := ss.listener.Accept()
		if err != nil {
			core.Debugf(ss.Name+"quit(err=%v)\n", err)
			break
		}
		core.Debugln(ss.Name + "connected to " + conn.RemoteAddr().String())
		conn.Close()
		core.Debugln(ss.Name + "disconnected " + conn.RemoteAddr().String())
	}
}

func (ss *SilentServer) Quit() {
	ss.listener.Close()
}
