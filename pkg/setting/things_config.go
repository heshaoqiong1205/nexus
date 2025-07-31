package setting

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type mqtt struct {
	Url              string `json:"url" ini:"url"`
	Keepalive        int    `json:"keep_alive" ini:"keepalive"`
	StateTopic       string `json:"state_topic" ini:"state_topic"`
	StateDesireTopic string `json:"state_desire_topic" ini:"state_desired_topic"`
	MetricsTopic     string `json:"metrics_topic" ini:"metrics_topic"`
	EventTopic       string `json:"event_topic" ini:"event_topic"`
	ServiceTopic     string `json:"service_topic" ini:"service_topic"`
}

type event struct {
	Interval    int    `json:"interval" ini:"interval"`
	Medium      string `json:"medium" ini:"medium"`
	MaxDuration int    `json:"max_duration" ini:"max_duration"`
}

type rtc struct {
	MaxDuration int `json:"max_duration" ini:"max_duration"`
	MaxRate     int `json:"max_rate" ini:"max_rate"`
}

type ThingsConfig struct {
	Mqtt mqtt `json:"mqtt"`
	Event event `json:"event"`
	Rtc rtc `json:"rtc"`
}

var thingsConfig ThingsConfig

func init() {
	thingsConfig = initThingsConfig()
	log.Printf("Things configuration loaded: %+v", thingsConfig)
}

func GetThingsConfig() *ThingsConfig {
	// This function should return the configuration for things.
	// For now, we return an empty struct as a placeholder.
	return &thingsConfig
}

func initThingsConfig() ThingsConfig {
    // Define possible config file paths
    configPaths := []string{
        "conf/things.ini",
        "../conf/things.ini",
        "../../conf/things.ini",
    }

    var cfg *ini.File
    var err error

    // Try to load the config file from one of the paths
    for _, path := range configPaths {
        if _, err = os.Stat(path); err == nil {
            cfg, err = ini.Load(path)
            if err == nil {
                log.Printf("Successfully loaded config file: %s", path)
                break
            }
        }
    }

    // If config file not found or couldn't be loaded, use defaults
    if err != nil {
        log.Printf("Warning: Could not load things.ini: %v. Using default configuration.", err)
        return getDefaultThingsConfig()
    }

    // Parse the INI file manually
    config := ThingsConfig{}
    
    // Parse MQTT section
    mqttSection := cfg.Section("mqtt")
    config.Mqtt.Url = mqttSection.Key("url").String()
    config.Mqtt.Keepalive = mqttSection.Key("keepalive").MustInt(60)
    config.Mqtt.StateTopic = mqttSection.Key("state_topic").String()
    config.Mqtt.StateDesireTopic = mqttSection.Key("state_desired_topic").String()
    config.Mqtt.MetricsTopic = mqttSection.Key("metrics_topic").String()
    config.Mqtt.EventTopic = mqttSection.Key("event_topic").String()
    config.Mqtt.ServiceTopic = mqttSection.Key("service_topic").String()
    
    // Parse Event section
    eventSection := cfg.Section("event")
    config.Event.Interval = eventSection.Key("interval").MustInt(60)
    config.Event.Medium = eventSection.Key("medium").String()
    config.Event.MaxDuration = eventSection.Key("max_duration").MustInt(3600)
    
    // Parse RTC section
    rtcSection := cfg.Section("rtc")
    config.Rtc.MaxDuration = rtcSection.Key("max_duration").MustInt(3600)
    config.Rtc.MaxRate = rtcSection.Key("max_rate").MustInt(1024)

    return config
}

// getDefaultThingsConfig returns default configuration values
func getDefaultThingsConfig() ThingsConfig {
    return ThingsConfig{
        Mqtt: mqtt{
            Url:            "tcp://localhost:1883",
            Keepalive:      60,
            StateTopic:     "things/${thingId}/state",
            StateDesireTopic: "things/${thingId}/state/desired",
            MetricsTopic:   "things/${thingId}/metrics",
            EventTopic:     "things/${thingId}/event",
            ServiceTopic:   "things/${thingId}/service",
        },
        Event: event{
            Interval:    60,
            Medium:      "mqtt",
            MaxDuration: 3600,
        },
        Rtc: rtc{
            MaxDuration: 3600,
            MaxRate:     1024,
        },
    }
}
