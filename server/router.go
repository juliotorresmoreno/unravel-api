package server

import "github.com/gin-gonic/gin"

func SetupServer() *gin.Engine {
	svr := gin.Default()

	svr.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return svr
}
