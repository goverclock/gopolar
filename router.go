package gopolar

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
)

// set up unix domain socket and signal handler for clean up
func (tm *TunnelManager) setupSock() {
	os.Remove("/tmp/gopolar.sock")
	sock, err := net.Listen("unix", "/tmp/gopolar.sock")
	if err != nil {
		log.Fatal(err)
	}
	tm.sock = sock

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove("/tmp/gopolar.sock")
		os.Exit(1)
	}()
}

func (tm *TunnelManager) setupRouter() {
	router := gin.Default()

	router.GET("/tunnels/list", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
				Tunnels []Tunnel `json:"tunnels"`
			} `json:"data"`
		}
		response.Success = true
		defer ctx.JSON(http.StatusOK, response)
		response.Data.Tunnels = tm.getTunnels() // never errors
	})

	router.POST("/tunnels/create", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
				ID uint64 `json:"id"`
			} `json:"data"`
		}
		response.Success = true
		defer ctx.JSON(http.StatusOK, response)
		request := CreateTunnelBody{}
		ctx.Bind(&request)

		newTunnel := Tunnel{
			// ID:     ,
			Name:   request.Name,
			Enable: false,
			Source: request.Source,
			Dest:   request.Dest,
		}
		newTunnelID, err := tm.addTunnel(newTunnel)
		if err != nil {
			response.Success = false
			response.ErrMsg = fmt.Sprint(err)
		} else {
			response.Data.ID = newTunnelID
		}
	})

	router.POST("/tunnels/edit/:id", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
			} `json:"data"`
		}
		response.Success = true
		defer ctx.JSON(http.StatusOK, response)
		request := EditTunnelBody{}
		ctx.Bind(&request)
		// parse params manually, due to possible gin issue on post params
		reqUrl := ctx.Request.URL.String()
		idStr := reqUrl[len("/tunnels/edit/"):]
		// idStr := ctx.Param("id")	// somehow buggy

		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			response.Success = false
			response.ErrMsg = fmt.Sprint(err)
			return
		}
		err = tm.changeTunnel(id, request.NewName, request.NewSource, request.NewDest)
		if err != nil {
			response.Success = false
			response.ErrMsg = fmt.Sprint(err)
		}
	})

	router.POST("/tunnels/toggle/:id", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
			} `json:"data"`
		}
		response.Success = true
		defer ctx.JSON(http.StatusOK, response)
		reqUrl := ctx.Request.URL.String()
		idStr := reqUrl[len("/tunnels/toggle/"):]
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			response.Success = false
			response.ErrMsg = fmt.Sprint(err)
			return
		}

		err = tm.toggleTunnel(id)
		if err != nil {
			response.Success = false
			response.ErrMsg = fmt.Sprint(err)
		}
	})

	router.DELETE("/tunnels/delete/:id", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
			} `json:"data"`
		}
		response.Success = true
		defer ctx.JSON(http.StatusOK, response)
		reqUrl := ctx.Request.URL.String()
		idStr := reqUrl[len("/tunnels/delete/"):]
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			response.Success = false
			response.ErrMsg = fmt.Sprint(err)
			return
		}

		err = tm.removeTunnel(id)
		if err != nil {
			response.Success = false
			response.ErrMsg = fmt.Sprint(err)
		}
	})

	router.GET("/about", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
				About AboutInfo `json:"about"`
			} `json:"data"`
		}
		response.Success = true
		defer ctx.JSON(http.StatusOK, response)

		response.Data.About = AboutInfo{
			Version: "1.0.0",
		}
	})

	tm.router = router
}
