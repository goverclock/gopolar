package testutil

import (
	"fmt"
	"net"
	"testing"
)

type EchoClient struct {
	name string
	port string
	conn net.Conn
	t    *testing.T
}

func NewEchoClient(t *testing.T, port uint64) *EchoClient {
	p := ":" + fmt.Sprint(port)
	ret := &EchoClient{
		name: "[client" + p + "] ",
		port: p,
		t:    t,
	}
	return ret
}

func (ec *EchoClient) Connect() error {
	conn, err := net.Dial("tcp", ec.port)
	ec.conn = conn
	if err != nil {
		ec.t.Log(ec.name + "fail to connect to " + conn.RemoteAddr().String())
	} else {
		ec.t.Log(ec.name + "connected to " + conn.RemoteAddr().String())
	}
	return err
}

func (ec *EchoClient) Disconnect() {
	ec.conn.Close()
	ec.conn = nil
	ec.t.Log(ec.name + "disconnected" + ec.port)
}

func (ec *EchoClient) Send(msg string) error {
	if msg[len(msg)-1] != '\n' {
		panic(ec.name + "trying to Send() without new line")
	}
	if ec.conn == nil {
		panic(ec.name + "trying to Send() without connection")
	}
	nw, err := ec.conn.Write([]byte(msg))
	if nw != len(msg) {
		return fmt.Errorf("partial write")
	}
	ec.t.Log(ec.name + "write " + msg)
	return err
}

func (ec *EchoClient) Recv() string {
	if ec.conn == nil {
		panic(ec.name + "trying to Recv() without connection")
	}
	buf := make([]byte, 1024*32)
	nr, err := ec.conn.Read(buf)
	if err != nil {
		panic(ec.name + err.Error())
	}
	reply := string(buf[:nr])
	ec.t.Log(ec.name + "read " + string(reply))
	return string(reply)

}
