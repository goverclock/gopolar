package main

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func handler() {
	sock, err := net.Listen("unix", "/tmp/gopolar.sock")
	check(err)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove("/tmp/gopolar.sock")
		os.Exit(1)
	}()

	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hhh")
	})
	router.RunListener(sock)
}
