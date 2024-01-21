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
	Name    string
	Port    uint64
	TotRecv uint64
	TotSend uint64

	conn       net.Conn // for send only, do not read conn, read lineReader
	lineReader *bufio.Reader
}

func NewEchoClient(port uint64) *EchoClient {
	p := ":" + fmt.Sprint(port)
	ret := &EchoClient{
		Name: "[client" + p + "] ",
		Port: port,
	}
	return ret
}

func (ec *EchoClient) Connect() error {
	conn, err := net.Dial("tcp", ":"+fmt.Sprint(ec.Port))
	ec.conn = conn
	ec.lineReader = bufio.NewReader(ec.conn)
	if err != nil {
		core.Debugln(ec.Name + "fail to connect to " + conn.RemoteAddr().String())
	} else {
		core.Debugln(ec.Name + "connected to " + conn.RemoteAddr().String())
	}
	return err
}

func (ec *EchoClient) Disconnect() {
	ec.conn.Close()
	ec.conn = nil
	core.Debugf(ec.Name+"disconnected %v", ec.Port)
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
		panic(ec.Name + "trying to Send() without new line")
	}
	if ec.conn == nil {
		panic(ec.Name + "trying to Send() without connection")
	}
	nw, err := ec.conn.Write([]byte(msg))
	if nw != len(msg) {
		return fmt.Errorf("partial write: %v", err)
	}
	core.Debugln(ec.Name + "send: " + msg)
	return err
}

// read until new line or EOF
func (ec *EchoClient) Recv() string {
	if ec.lineReader == nil {
		panic(ec.Name + "trying to Recv() without connection")
	}

	reply := []byte{}
	// core.Debugln(ec.name + "reading")
	bytes, err := ec.lineReader.ReadBytes(byte('\n'))
	core.Debugf(ec.Name+"read %v bytes\n", len(bytes))
	if err == nil || errors.Is(err, io.EOF) {
		reply = append(reply, bytes...)
	} else {
		panic(ec.Name + err.Error())
	}

	s := string(reply)
	core.Debugln(ec.Name + "recv: " + s[:len(s)-1] + "\\n")
	ec.TotRecv += uint64(len(s))
	return s
}
