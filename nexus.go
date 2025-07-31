package main

import (
	"nexus/http_server"
	"nexus/models"
	"nexus/pkg/setting"
)

func main() {
	setting.Setup()
	models.Setup(setting.DatabaseSetting)
	if http_server.Run(setting.ServerSetting) != nil {
		panic("Failed to start HTTP server")
	}
}
