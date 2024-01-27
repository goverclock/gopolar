package core

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func (tm *TunnelManager) setupRouter() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New() // use gin.Default() for http log, gin.New() to omit
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
	}))
	router.Use(static.Serve("/", static.LocalFile("./gpwebui", false)))

	router.GET("/tunnels/list", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
				Tunnels []Tunnel `json:"tunnels"`
			} `json:"data"`
		}
		response.Success = true
		response.Data.Tunnels = tm.GetTunnels()
		ctx.JSON(http.StatusOK, response)
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
		request := CreateTunnelBody{}
		ctx.Bind(&request)

		newTunnel := Tunnel{
			// ID:     ,
			Name:   request.Name,
			Enable: true, // new tunnels are enabled by default
			Source: request.Source,
			Dest:   request.Dest,
		}
		newTunnelID, err := tm.AddTunnel(newTunnel)
		if err != nil {
			response.Success = false
			response.ErrMsg = fmt.Sprint(err)
		} else {
			response.Data.ID = newTunnelID
		}
		ctx.JSON(http.StatusOK, response)
	})

	router.POST("/tunnels/edit/:id", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
			} `json:"data"`
		}
		response.Success = true
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
			ctx.JSON(http.StatusOK, response)
			return
		}
		err = tm.ChangeTunnel(id, request.NewName, request.NewSource, request.NewDest)
		if err != nil {
			response.Success = false
			response.ErrMsg = fmt.Sprint(err)
		}
		ctx.JSON(http.StatusOK, response)
	})

	router.POST("/tunnels/toggle/:id", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
			} `json:"data"`
		}
		response.Success = true
		reqUrl := ctx.Request.URL.String()
		idStr := reqUrl[len("/tunnels/toggle/"):]
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			response.Success = false
			response.ErrMsg = fmt.Sprint(err)
			ctx.JSON(http.StatusOK, response)
			return
		}

		err = tm.ToggleTunnel(id)
		if err != nil {
			response.Success = false
			response.ErrMsg = fmt.Sprint(err)
		}
		ctx.JSON(http.StatusOK, response)
	})

	router.DELETE("/tunnels/delete/:id", func(ctx *gin.Context) {
		var response struct {
			Success bool   `json:"success"`
			ErrMsg  string `json:"err_msg"`
			Data    struct {
			} `json:"data"`
		}
		response.Success = true
		reqUrl := ctx.Request.URL.String()
		idStr := reqUrl[len("/tunnels/delete/"):]
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			response.Success = false
			response.ErrMsg = fmt.Sprint(err)
			ctx.JSON(http.StatusOK, response)
			return
		}

		err = tm.RemoveTunnel(id)
		if err != nil {
			response.Success = false
			response.ErrMsg = fmt.Sprint(err)
		}
		ctx.JSON(http.StatusOK, response)
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

		response.Data.About = AboutInfo{
			Version: "1.0.0",
		}
		ctx.JSON(http.StatusOK, response)
	})

	tm.router = router
}
