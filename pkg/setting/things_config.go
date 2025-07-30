package setting

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type mqtt struct {
	Url 	string `json:"url"`
	Keepalive int    `json:"keep_=alive"`
	StateTopic string `json:"state_topic"`
	StateDesireTopic string `json:"state_desire_topic"`
	MetricsTopic string `json:"metrics_topic"`
	EventTopic string `json:"event_topic"`
	ServiceTopic string `json:"service_topic"`
}

type event struct {
	Interval   int    `json:"interval"`
	Medium    string `json:"medium"`
	MaxDuration int    `json:"max_duration"`
}

type rtc struct {
	MaxDuration int `json:"max_duration"`
	MaxRate int `json:"max_rate"`
}

type ThingsConfig struct {
	Mqtt mqtt `json:"mqtt"`
	Event event `json:"event"`
	Rtc rtc `json:"rtc"`
}

var thingsConfig ThingsConfig

func init() {
	thingsConfig = initThingsConfig()
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
                break
            }
        }
    }

    // If config file not found or couldn't be loaded, use defaults
    if err != nil {
        log.Printf("Warning: Could not load things.ini: %v. Using default configuration.", err)
        return getDefaultThingsConfig()
    }

    // Parse the INI file into ThingsConfig struct
    var config ThingsConfig
    err = cfg.MapTo(&config)
    if err != nil {
        log.Printf("Error parsing things.ini: %v. Using default configuration.", err)
        return getDefaultThingsConfig()
    }

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
