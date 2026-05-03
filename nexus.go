package main

import (
	"nexus/http_server"
	"nexus/models"
	"nexus/pkg/setting"
	"nexus/services/mqtt"
	"nexus/services/things"
)

func main() {
	setting.Setup()
	models.Setup(setting.DatabaseSetting)

	if setting.GetThingsConfig().Mqtt.Enabled {
		thingService := things.NewThingService()
		mqttRuntime, err := mqtt.NewRuntime(setting.GetThingsConfig(), mqtt.NewServiceDispatcher(thingService))
		if err != nil {
			panic("Failed to initialize MQTT runtime: " + err.Error())
		}
		if err := mqttRuntime.Start(); err != nil {
			panic("Failed to start MQTT runtime: " + err.Error())
		}
		defer mqttRuntime.Stop()
		things.SetDefaultRPCTransport(mqttRuntime)
	}

	if http_server.Run(setting.ServerSetting) != nil {
		panic("Failed to start HTTP server")
	}
}
