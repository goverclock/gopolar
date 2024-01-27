package tui

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/goverclock/gopolar/internal/core"

	"github.com/gin-gonic/gin"
)

func setupMockRouter(e *gin.Engine) {
	log.Println("setup mock router")

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

	tunnels := []core.Tunnel{
		{
			ID:     1,
			Name:   "first tunnel",
			Enable: false,
			Source: "localhost:2345",
			Dest:   "233.168.10.1:5678",
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
				Tunnels []core.Tunnel `json:"tunnels"`
			} `json:"data"`
		}
		response.Success = true
		response.Data.Tunnels = append(response.Data.Tunnels, tunnels...)
		ctx.JSON(http.StatusOK, response)
	})

	mock_router.POST("/tunnels/create", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
				ID uint64 `json:"id"`
			} `json:"data"`
		}
		request := core.CreateTunnelBody{}
		ctx.Bind(&request)
		log.Printf("%#v", request)
		newTunnel := core.Tunnel{
			ID:     rand.Uint64() % 100,
			Name:   request.Name,
			Enable: false,
			Source: request.Source,
			Dest:   request.Dest,
		}
		tunnels = append(tunnels, newTunnel)
		response.Success = true
		response.Data.ID = newTunnel.ID
		ctx.JSON(http.StatusOK, response)
	})

	mock_router.POST("/tunnels/edit/:id", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
			} `json:"data"`
		}
		request := core.EditTunnelBody{}
		ctx.Bind(&request)
		log.Printf("%#v", request)
		// parse params manually, due to possible gin issue on post params
		reqUrl := ctx.Request.URL.String()
		idStr := reqUrl[len("/tunnels/edit/"):]
		// idStr := ctx.Param("id")	// somehow buggy
		_, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Println(err)
		}
		for i, t := range tunnels {
			if fmt.Sprint(t.ID) == idStr {
				tunnels[i].Name, tunnels[i].Source, tunnels[i].Dest = request.NewName, request.NewSource, request.NewDest
			}
		}
		response.Success = true
		ctx.JSON(http.StatusOK, response)
	})

	mock_router.POST("/tunnels/toggle/:id", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
			} `json:"data"`
		}
		reqUrl := ctx.Request.URL.String()
		idStr := reqUrl[len("/tunnels/toggle/"):]
		_, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Println(err)
		}
		for i, t := range tunnels {
			if fmt.Sprint(t.ID) == idStr {
				tunnels[i].Enable = !tunnels[i].Enable
			}
		}
		response.Success = true
		ctx.JSON(http.StatusOK, response)
	})

	mock_router.DELETE("/tunnels/delete/:id", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
			} `json:"data"`
		}
		reqUrl := ctx.Request.URL.String()
		idStr := reqUrl[len("/tunnels/delete/"):]
		id, err := strconv.ParseUint(idStr, 10, 64)
		for i, t := range tunnels {
			if t.ID == id {
				tunnels = append(tunnels[:i], tunnels[i+1:]...)
				break
			}
		}
		if err != nil {
			log.Println(err)
		}
		response.Success = true
		ctx.JSON(http.StatusOK, response)
	})

	target := core.AboutInfo{
		Version: "0.0.1",
	}
	mock_router.GET("/about", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
				About core.AboutInfo `json:"about"`
			} `json:"data"`
		}
		response.Data.About = target
		response.Success = true
		ctx.JSON(http.StatusOK, response)
	})
}
