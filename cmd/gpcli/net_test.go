package main

import (
	"gopolar"
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

func TestMain(m *testing.M) {
	// init
	log.SetPrefix("[net_test]")
	log.Println("listening on unix domain socket")
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

	exitCode := m.Run()
	cleanup()
	os.Exit(exitCode)
}

func cleanup() {
	os.Remove("/tmp/gopolar.sock")
}

// test communication functionality
func TestGetGood(t *testing.T) {
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

func TestGetBad(t *testing.T) {
	assert := assert.New(t)

	m := end.GET("/notexist")
	t.Log(m)
	assert.Equal(0, len(m))
}

func TestGetTunnelsList(t *testing.T) {
	assert := assert.New(t)

	expectList := []gopolar.Tunnel{
		{
			ID:     1,
			Name:   "first tunnel",
			Enable: false,
			Source: "localhost:2345",
			Dest:   "192.168.10.1:4567",
		},
		{
			ID:     2,
			Name:   "second tunnel",
			Enable: false,
			Source: "localhost:3333",
			Dest:   "192.168.10.1:4567",
		},
		{
			ID:     3,
			Name:   "hahaha this is 3",
			Enable: true,
			Source: "localhost:2789",
			Dest:   "localhost:2333",
		},
	}
	mock_router.GET("/tunnels/list", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
				Tunnels []gopolar.Tunnel `json:"tunnels"`
			} `json:"data"`
		}
		response.Success = true
		response.Data.Tunnels = append(response.Data.Tunnels, expectList...)
		ctx.JSON(http.StatusOK, response)
	})

	retList, err := end.GetTunnelsList()
	assert.Equal(nil, err)
	assert.Equal(expectList, retList)
}
