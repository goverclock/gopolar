package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

var srcPort int = 8001
var dstPort int = 8002

func main() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(srcPort))
	if err != nil {
		panic("failed to listen" + err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("fail to accept", err)
			continue
		}
		copyConn(conn)
	}
}

func copyConn(src net.Conn) {
	dst, err := net.Dial("tcp", ":"+strconv.Itoa(dstPort))
	if err != nil {
		panic("fail to dial" + err.Error())
	}

	done := make(chan struct{})

	// src -> dst
	go func() {
		defer src.Close()
		defer dst.Close()
		buf := make([]byte, 1024)
		n, _ := src.Read(buf)
		fmt.Printf("%d bytes forwarded: %v\n", n, buf[:n])
		dst.Write(buf[:n])
		// io.Copy(dst, src)
		done <- struct{}{}
	}()
	// dst->src
	go func() {
		defer src.Close()
		defer dst.Close()
		io.Copy(src, dst)
		done <- struct{}{}
	}()

	<-done
	<-done
}
