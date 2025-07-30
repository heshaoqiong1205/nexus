package http_server

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"nexus/pkg/setting"
	"nexus/services/storage"
	"nexus/services/things"
	"strings"
	"time"

	external_storage "nexus/external/storage"

	"github.com/gin-gonic/gin"
)

// DeviceActivationRequest represents the device activation request structure

// DeviceActivationResponse represents the device activation response structure
const (
	BadRequest = 1

)

var thingsService *things.ThingService
var cloudStorageService *storage.CloudStorageService

func init() {
	thingsService = things.NewThingService()
	cloudStorageService = storage.NewCloudStorageService()
}

type Response struct {
    Success     bool   `json:"success"`
    Result      interface{} `json:"result,omitempty"`
}

type ActiveResult struct {
	Device *things.IoTDevice `json:"device"`
	Config *setting.ThingsConfig `json:"config"`
}

func Run() error {
    route := gin.Default()

	route.Use(DeviceAuthMiddleware())

    route.GET("/things/device/:id/state", func(c *gin.Context) {
        c.JSON(200, gin.H{"id": "device.ID"})
    })

    // Device activation
    route.POST("/things/device/activate", activateDevice)

	// Device deactivate
	route.POST("/things/device/deactivate", deactivateDevice)

	// Device config
	route.GET("/things/device/config", getDeviceConfig)
	// Device get config

	route.GET("/things/device/credentials", generateGetStorageCredentials)

    return route.Run(":8080")
}

func DeviceAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip auth for device activation endpoint
        if c.FullPath() == "/things/device/activate" {
            c.Next()
            return
        }

        // Check if path starts with "/things/device/"
        if !strings.HasPrefix(c.Request.URL.Path, "/things/device/") {
            c.Next()
            return
        }

        // Get required headers
        deviceID := c.GetHeader("device-id")
        reqTime := c.GetHeader("time")
        headerDigest := c.GetHeader("digest")

        if deviceID == "" || reqTime == "" || headerDigest == "" {
            c.JSON(http.StatusUnauthorized, Response{
                Success: false,
                Result:  "Missing required authentication headers",
            })
            c.Abort()
            return
        }

        // Get device secret key
        device, err := thingsService.GetDevice(deviceID)
        if err != nil {
            c.JSON(http.StatusUnauthorized, Response{
                Success: false,
                Result:  "Device authentication failed",
            })
            c.Abort()
            return
        }

        // Check request time to prevent replay attacks
        timestamp, err := time.Parse(time.RFC3339, reqTime)
        if err != nil || time.Since(timestamp).Minutes() > 5 {
            c.JSON(http.StatusUnauthorized, Response{
                Success: false,
                Result:  "Request expired or invalid timestamp",
            })
            c.Abort()
            return
        }

        // Read request body
        bodyBytes, err := io.ReadAll(c.Request.Body)
        if err != nil {
            c.JSON(http.StatusInternalServerError, Response{
                Success: false,
                Result:  "Error reading request body",
            })
            c.Abort()
            return
        }

        // Restore the request body for next handlers
        c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

        // Calculate body digest
        h := sha256.New()
        h.Write(bodyBytes)
        h.Write([]byte(device.SecretKey))
        calculatedDigest := hex.EncodeToString(h.Sum(nil))

        // Compare digests
        if headerDigest != calculatedDigest {
            c.JSON(http.StatusUnauthorized, Response{
                Success: false,
                Result:  "Invalid request signature",
            })
            c.Abort()
            return
        }

        // Authentication successful
        c.Next()
    }
}


// activateDevice handles the device activation request
func activateDevice(c *gin.Context) {
    var req things.ActiveRequest

    // Bind and validate request parameters
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusOK, Response{
            Success: false,
            Result: err.Error(),
        })
        return
    }


    // Simulate successful activation
    device, err := thingsService.Active(&req)
    if err != nil {
        c.JSON(http.StatusOK, Response{
            Success: false,
            Result: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Success:     true,
        Result:      device,
    })
}

func deactivateDevice(c *gin.Context) {
	deviceID := c.GetHeader("device-id")

	// Simulate successful deactivation
	err := thingsService.Deactivate(deviceID)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Success: false,
			Result: err.Error(),
		})
		return
	}

	result := ActiveResult{
		Device: nil, // Device is deactivated, so we return nil
		Config: setting.GetThingsConfig(),
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Result:  result,
	})
}



func getDeviceConfig(c *gin.Context) {
	config := setting.GetThingsConfig()
	c.JSON(http.StatusOK, Response{
		Success: true,
		Result:  config,
	})
}

func generateGetStorageCredentials(c *gin.Context) {
	deviceId := c.GetHeader("device-id")

	// Get storage credentials
	actions := []external_storage.Action{external_storage.GetObject, external_storage.PutObject}
	credentials, err := cloudStorageService.GetStorageCredentials(actions, deviceId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Success: false,
			Result:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Result:  credentials,
	})
}
