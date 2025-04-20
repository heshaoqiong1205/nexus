package main

import (
	"nexus/http_server"
	"nexus/pkg/setting"
)

func main() {
	setting.Setup()
	http_server.Run()
}
