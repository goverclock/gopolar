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

	m, err := end.GET("/hello")
	assert.Equal(nil, err)
	t.Log(m)
	assert.Equal(1, len(m))
	assert.Equal("world", m["hello"])
}

func TestGetBad(t *testing.T) {
	assert := assert.New(t)

	m, err := end.GET("/notexist")
	assert.NotEqual(nil, err)
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

func TestCreateTunnel(t *testing.T) {
	assert := assert.New(t)

	name := "created me"
	source := "localhost:3456"
	dest := "localhost:4567"
	mock_router.POST("/tunnels/create", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
			} `json:"data"`
		}
		var request struct {
			Name   string `json:"name"`
			Source string `json:"source"`
			Dest   string `json:"dest"`
		}
		ctx.Bind(&request)
		assert.Equal(name, request.Name)
		assert.Equal(source, request.Source)
		assert.Equal(dest, request.Dest)
		response.Success = true
		ctx.JSON(http.StatusOK, response)
	})

	err := end.CreateTunnel(name, source, dest)
	log.Println(err)
	assert.Equal(nil, err)
}
