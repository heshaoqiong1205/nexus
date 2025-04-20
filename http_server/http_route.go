package http_server

import (
	"github.com/gin-gonic/gin"
)

func Run() {
	route := gin.Default()
	route.GET("/things/device/:id/state", func(c *gin.Context) {
		c.JSON(200, gin.H{"id": "device.ID"})
	})
	route.Run(":8080")
}
