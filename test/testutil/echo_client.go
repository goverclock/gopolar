package testutil

import (
	"bufio"
	"errors"
	"fmt"
	"gopolar/internal/core"
	"io"
	"net"
	"time"
)

type EchoClient struct {
	name       string
	port       string
	conn       net.Conn // for send only, do not read conn, read lineReader
	lineReader *bufio.Reader
}

func NewEchoClient(port uint64) *EchoClient {
	p := ":" + fmt.Sprint(port)
	ret := &EchoClient{
		name: "[client" + p + "] ",
		port: p,
	}
	return ret
}

func (ec *EchoClient) Connect() error {
	conn, err := net.Dial("tcp", ec.port)
	ec.conn = conn
	ec.lineReader = bufio.NewReader(ec.conn)
	if err != nil {
		core.Debugln(ec.name + "fail to connect to " + conn.RemoteAddr().String())
	} else {
		core.Debugln(ec.name + "connected to " + conn.RemoteAddr().String())
	}
	return err
}

func (ec *EchoClient) Disconnect() {
	ec.conn.Close()
	ec.conn = nil
	core.Debugln(ec.name + "disconnected" + ec.port)
}

// this would try to read a byte,
// do not Recv() after this
func (ec *EchoClient) IsConnected() bool {
	if ec.conn == nil {
		return false
	}
	one := make([]byte, 1)
	ec.conn.SetReadDeadline(time.Now().Add(time.Millisecond))
	if _, err := ec.conn.Read(one); err == io.EOF {
		ec.conn.Close()
		ec.conn = nil
		return false
	}
	return true
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
	core.Debugln(ec.name + "send: " + msg)
	return err
}

// read until new line or EOF
func (ec *EchoClient) Recv() string {
	if ec.lineReader == nil {
		panic(ec.name + "trying to Recv() without connection")
	}

	reply := []byte{}
	// core.Debugln(ec.name + "reading")
	bytes, err := ec.lineReader.ReadBytes(byte('\n'))
	core.Debugf(ec.name+"read %v bytes\n", len(bytes))
	if err == nil || errors.Is(err, io.EOF) {
		reply = append(reply, bytes...)
	} else {
		panic(ec.name + err.Error())
	}

	s := string(reply)
	core.Debugln(ec.name + "recv: " + s[:len(s)-1] + "\\n")
	return string(reply)
}
