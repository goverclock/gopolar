package main

import (
	"gopolar"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var end *CLIEnd
var mock_router *gin.Engine

func TestMain(m *testing.M) {
	// init
	log.SetPrefix("[gpcli_test]")
	log.Println("listening on unix domain socket")
	log.SetFlags(0)
	os.Remove("/tmp/gopolar.sock")
	isMock := os.Getenv("MOCK")

	end = NewCLIEnd()
	// mock core
	sock, err := net.Listen("unix", "/tmp/gopolar.sock")
	if err != nil {
		log.Fatal(err)
	}

	mock_router = gin.Default()
	exitCode := 0
	if isMock == "1" {
		setupMockRouter(mock_router)
		mock_router.RunListener(sock)
	} else {
		go mock_router.RunListener(sock)
		exitCode = m.Run()
	}

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
	createdTunnelID := uint64(23423)
	mock_router.POST("/tunnels/create", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
				ID uint64 `json:"id"`
			} `json:"data"`
		}
		request := gopolar.CreateTunnelBody{}
		ctx.Bind(&request)
		assert.Equal(name, request.Name)
		assert.Equal(source, request.Source)
		assert.Equal(dest, request.Dest)
		response.Success = true
		response.Data.ID = createdTunnelID
		ctx.JSON(http.StatusOK, response)
	})

	id, err := end.CreateTunnel(name, source, dest)
	assert.Equal(nil, err)
	assert.Equal(createdTunnelID, id)
}

func TestEditTunnel(t *testing.T) {
	assert := assert.New(t)

	targetID := int64(12345)
	newName := "new created me"
	newSource := "newhahah:3456"
	newDest := "newDest:4567"
	mock_router.POST("/tunnels/edit/:id", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
			} `json:"data"`
		}
		request := gopolar.EditTunnelBody{}
		ctx.Bind(&request)
		// parse params manually, due to possible gin issue on post params
		reqUrl := ctx.Request.URL.String()
		idStr := reqUrl[len("/tunnels/edit/"):]
		// idStr := ctx.Param("id")	// somehow buggy
		// log.Println("idStr=", idStr)
		// log.Println("request.url=", ctx.Request.URL)
		// log.Println("params=", ctx.Params)
		recvID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(targetID, recvID)
		assert.Equal(newName, request.NewName)
		assert.Equal(newSource, request.NewSource)
		assert.Equal(newDest, request.NewDest)
		response.Success = true
		ctx.JSON(http.StatusOK, response)
	})

	err := end.EditTunnel(targetID, newName, newSource, newDest)
	assert.Equal(nil, err)
}

func TestToggleTunnel(t *testing.T) {
	assert := assert.New(t)

	targetID := int64(56785)
	mock_router.POST("/tunnels/toggle/:id", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
			} `json:"data"`
		}
		reqUrl := ctx.Request.URL.String()
		idStr := reqUrl[len("/tunnels/toggle/"):]
		recvID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(targetID, recvID)
		response.Success = true
		ctx.JSON(http.StatusOK, response)
	})
	err := end.ToggleTunnel(targetID)
	assert.Equal(nil, err)
}

func TestDeleteTunnel(t *testing.T) {
	assert := assert.New(t)

	targetID := int64(78967)
	mock_router.DELETE("/tunnels/delete/:id", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
			} `json:"data"`
		}
		reqUrl := ctx.Request.URL.String()
		idStr := reqUrl[len("/tunnels/delete/"):]
		recvID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(targetID, recvID)
		response.Success = true
		ctx.JSON(http.StatusOK, response)
	})
	err := end.DeleteTunnel(targetID)
	assert.Equal(nil, err)
}

func TestGetAboutInfo(t *testing.T) {
	assert := assert.New(t)

	target := gopolar.AboutInfo{
		Version: "0.0.1",
	}
	mock_router.GET("/about", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
				About gopolar.AboutInfo `json:"about"`
			} `json:"data"`
		}
		response.Data.About = target
		response.Success = true
		ctx.JSON(http.StatusOK, response)
	})
	ret, err := end.GetAboutInfo()
	assert.Equal(nil, err)
	assert.Equal(target, ret)
}
