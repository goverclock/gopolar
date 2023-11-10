package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	go func() {
		listener, err := net.Listen("tcp", ":8002")
		if err != nil {
			panic("failed to listen" + err.Error())
		}
		conn, err := listener.Accept()
		if err != nil {
			panic("fail to accept" + err.Error())
		}
        defer conn.Close()

		reply := make([]byte, 1024)
		n, err := conn.Read(reply)
		if err != nil {
			panic("fail to read" + err.Error())
		}
		fmt.Printf("read %d bytes from 8002: %v", n, string(reply))
	}()

	time.Sleep(1 * time.Second)
	connIn, err := net.Dial("tcp", ":8001")
	if err != nil {
		panic("fail to dial in" + err.Error())
	}
	defer connIn.Close()

    message := "hello"
	n, err := connIn.Write([]byte(message))
	if err != nil {
		panic("fail to write to connIn" + err.Error())
	}
	fmt.Printf("wrote %d bytes to 8001: %v\n", n, message)

	time.Sleep(1 * time.Second)
}
