package server

import (
	"github.com/gin-gonic/gin"
	"github.com/juliotorresmoreno/unravel-api/server/events"
)

func SetupServer() *gin.Engine {
	svr := gin.Default()

	events.SetupRouter(svr.Group("/events"))

	svr.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return svr
}
