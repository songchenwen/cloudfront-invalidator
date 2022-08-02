package server

import (
	"github.com/gin-gonic/gin"
	"github.com/songchenwen/cloudfront-invalidator/config"
)

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	if config.IsDebug() {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.POST("/:distribution", handleInvalidate)
	r.GET("/:distribution", handleInvalidate)
	return r
}
