package main

import (
	"gopolar"
	"log"
	"net/http"
	"strconv"

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
		response.Success = true
		response.Data.ID = createdTunnelID
		ctx.JSON(http.StatusOK, response)
	})

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
		_, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Println(err)
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
		_, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Println(err)
		}
		response.Success = true
		ctx.JSON(http.StatusOK, response)
	})

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
}
