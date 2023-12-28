package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var end *CLIEnd
var mock_router *gin.Engine

func init() {
	log.SetFlags(0)
	os.Remove("/tmp/gopolar.sock")

	end = NewCLIEnd()
	// mock core
	sock, err := net.Listen("unix", "/tmp/gopolar.sock")
	if err != nil {
		log.Fatal(err)
	}

	mock_router = gin.Default()
	go mock_router.RunListener(sock)
	log.SetPrefix("[net_test]")
	log.Println("listening on unix domain socket")
}

func cleanup() {
	os.Remove("/tmp/gopolar.sock")
}

// test communication functionality
func TestGetHello(t *testing.T) {
	t.Cleanup(cleanup)
	assert := assert.New(t)
	mock_router.GET("/hello", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
				Hello string `json:"hello"`
			} `json:"data"`
		}
		response.Success = true
		response.Data.Hello = "world"
		ctx.JSON(http.StatusOK, response)
	})

	m := end.GET("/hello")
	t.Log(m)
	assert.Equal(1, len(m))
	assert.Equal("world", m["hello"])
}


