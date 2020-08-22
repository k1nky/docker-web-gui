package main

import (
	"dockerboard/controller"
	"net/http"

	"github.com/gin-contrib/static"

	"github.com/gin-gonic/gin"
)

func main() {
	Cfg.ParseArgs()
	r := gin.New()
	authMW := controller.NewAuthMiddleware("secret#")
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/api/container/stats"},
	}))
	r.NoRoute(NoRouteHandler)
	r.Use(static.Serve("/", static.LocalFile("./ui", true)))
	auth := r.Group("/api")
	auth.GET("/refresh_token", authMW.RefreshHandler)
	auth.POST("/login", authMW.LoginHandler)
	auth.Use(authMW.MiddlewareFunc())
	{
		auth.GET("/image/fetch", controller.GetImageHandler(Cfg.DockerHost))
		auth.GET("/container/fetch", controller.GetContainerHandler(Cfg.DockerHost))
		auth.OPTIONS("/container/fetch", controller.GetContainerHandler(Cfg.DockerHost))
		auth.GET("/container/logs", controller.GetLogsHandler(Cfg.DockerHost))
		auth.GET("/container/stats", controller.GetStatsHandler(Cfg.DockerHost))
		auth.GET("/container/command", controller.GetCommandHandler(Cfg.DockerHost))
	}
	r.RunTLS(Cfg.Listen, "cert.pem", "key.pem")
}

func NoRouteHandler(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, "/")
}
