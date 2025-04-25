package http_server

import (
	"github.com/gin-gonic/gin"
)

func Run() error {
	route := gin.Default()
	route.GET("/things/device/:id/state", func(c *gin.Context) {
		c.JSON(200, gin.H{"id": "device.ID"})
	})
	return route.Run(":8080")
}
