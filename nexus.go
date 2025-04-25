package main

import (
	"nexus/http_server"
	"nexus/pkg/setting"
)

func main() {
	setting.Setup()
	if http_server.Run() != nil {
		panic("Failed to start HTTP server")
	}
}
