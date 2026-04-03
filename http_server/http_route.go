package http_server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"nexus/pkg/setting"
	"nexus/services/auth"
	"nexus/services/storage"
	"nexus/services/things"
	"nexus/services/user"
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
var userService user.IUserService

func init() {
	thingsService = things.NewThingService()
	cloudStorageService = storage.NewCloudStorageService()
	userService = user.NewUserService()
}

type Response struct {
    Success   bool        `json:"success"`
	Timestamp int64   `json:"timestamp"`
    Result    interface{} `json:"result,omitempty"`
    Error     string      `json:"error,omitempty"`
}

// ErrorResponse creates a standardized error response
func ErrorResponse(message string) Response {
    return Response{
        Success: false,
		Timestamp: time.Now().UnixMilli(),
        Error:   message,
    }
}

// SuccessResponse creates a standardized success response
func SuccessResponse(data interface{}) Response {
    return Response{
        Success: true,
		Timestamp: time.Now().UnixMilli(),
        Result:  data,
    }
}

type ActiveResult struct {
	Device *things.IoTDevice `json:"device"`
	Config *setting.ThingsConfig `json:"config"`
}


// User

type SignInRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

func Run(config *setting.Server) error {
    // Set Gin mode based on configuration
    if config.RunMode == "release" {
        gin.SetMode(gin.ReleaseMode)
    } else {
        gin.SetMode(gin.DebugMode)
    }

    route := gin.Default()

    // Apply authentication middlewares
    route.Use(DeviceAuthMiddleware())
    route.Use(UserMiddleware())

	route.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok!")
	})
	route.POST("/cloud/account/sign_in", signIn)
	route.POST("/cloud/account/sign_up", signUp)

	route.GET("/cloud/user/details", userDetails)

	route.GET("/cloud/devices", deviceList)
    route.GET("/cloud/device/:device_id", deviceDetails)


    // Device activation
    route.POST("/things/activate", activateDevice)
    // Device deactivate
    route.POST("/things/deactivate", deactivateDevice)

    // Device config
    route.GET("/things/device/config", getDeviceConfig)
    route.GET("/things/device/credentials", generateGetStorageCredentials)

	// Device report event
	route.POST("/things/device/event", handleEvent)
	// Device report state
	route.POST("/things/device/state", handleState)
	// Get device desired state
	route.GET("/things/device/state/desired", fetchDesiredState)

    // Use configuration values for server setup
    server := &http.Server{
        Addr:         fmt.Sprintf(":%d", config.HttpPort),
        Handler:      route,
        ReadTimeout:  config.ReadTimeout,
        WriteTimeout: config.WriteTimeout,
    }

    log.Printf("Starting HTTP server on port %d in %s mode", config.HttpPort, config.RunMode)
    log.Printf("Server timeouts - Read: %v, Write: %v", config.ReadTimeout, config.WriteTimeout)

    return server.ListenAndServe()
}

func getDeviceID(c *gin.Context) string {
	return c.GetHeader("device-id")
}

func getDeviceDigest(c *gin.Context) string {
	return c.GetHeader("digest")
}

func getUserID(c *gin.Context) string {
	return c.GetString("user-id")
}

func getUserToken(c *gin.Context) string {
	// Get token from Authorization header (Bearer token)
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check if it starts with "Bearer "
	if len(authHeader) > 7 && strings.ToLower(authHeader[:7]) == "bearer " {
		return authHeader[7:] // Return token without "Bearer " prefix
	}

	// Also check for direct token header
	if token := c.GetHeader("token"); token != "" {
		return token
	}

	return ""
}

func DeviceAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check if path starts with "/things/device/"
        if !strings.HasPrefix(c.Request.URL.Path, "/things/device/") {
            c.Next()
            return
        }

		if c.Request.URL.Path == "/things/activate" {
			// Additional checks for device activation
			c.Next()
			return
		}

        // Get required headers
        deviceID := getDeviceID(c)
        deviceDigest := getDeviceDigest(c)

        if deviceID == "" || deviceDigest == "" {
            c.JSON(http.StatusUnauthorized, ErrorResponse("Missing required authentication headers"))
            c.Abort()
            return
        }

        // Get device secret key
        device, err := thingsService.GetDevice(deviceID)
        if err != nil {
            c.JSON(http.StatusUnauthorized, ErrorResponse("Device authentication failed"))
            c.Abort()
            return
        }


        // Calculate body digest
        h := hmac.New(sha256.New, []byte(device.SecretKey))
        h.Write([]byte(deviceID))
        calculatedDigest := base64.StdEncoding.EncodeToString(h.Sum(nil))

        // Compare digests
        if deviceDigest != calculatedDigest {
            c.JSON(http.StatusUnauthorized, ErrorResponse("Invalid request signature"))
            c.Abort()
            return
        }

        // Store device info in context for later use
        c.Set("device", device)
        c.Next()
    }
}

func UserMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip authentication for public endpoints
        path := c.Request.URL.Path
        if path == "/health" ||
           path == "/cloud/account/sign_in" ||
           path == "/cloud/account/sign_up" ||
           strings.HasPrefix(path, "/things/") {
            c.Next()
            return
        }

        // Only apply user authentication for /cloud/ paths
        if !strings.HasPrefix(path, "/cloud/") {
            c.Next()
            return
        }

        // Get token from header
        token := getUserToken(c)
        if token == "" {
            c.JSON(http.StatusUnauthorized, ErrorResponse("Missing or invalid authorization token"))
            c.Abort()
            return
        }

        // Validate token and get user info
        userID, err := auth.ValidateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, ErrorResponse("Invalid or expired token"))
            c.Abort()
            return
        }

        // Store user info in context for later use
        c.Set("user-id", userID)
        c.Set("user-token", token)
        c.Next()
    }
}

func deviceList(c *gin.Context) {

	var query things.DevicesQuery

	    // Bind and validate request parameters
    if err := c.ShouldBindJSON(&query); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request format: "+err.Error()))
        return
    }
	devices, err := thingsService.GetDevices(query)
	if err != nil {
		log.Printf("Error fetching device list: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("Failed to fetch device list"))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(devices))
}

func deviceDetails(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse("Missing device ID"))
		return
	}

	log.Printf("Fetching details for device ID: %s", deviceID)

	// Get device state
	device, err := thingsService.GetDevice(deviceID)
	if err != nil {
		log.Printf("Error fetching device %s: %v", deviceID, err)
		c.JSON(http.StatusNotFound, ErrorResponse("Device not found"))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(device))
}

// activateDevice handles the device activation request
func activateDevice(c *gin.Context) {
    var req things.ActiveRequest

    // Bind and validate request parameters
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request format: "+err.Error()))
        return
    }

	fmt.Printf("Activation request: %+v\n", req)
    // Activate device
    device, err := thingsService.Active(&req)
    if err != nil {
        log.Printf("Device activation failed: %v", err)
        c.JSON(http.StatusInternalServerError, ErrorResponse("Device activation failed: "+err.Error()))
        return
    }

    c.JSON(http.StatusOK, SuccessResponse(device))
}

func deactivateDevice(c *gin.Context) {
	deviceID := getDeviceID(c)
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse("Missing device-id header"))
		return
	}

	err := thingsService.Deactivate(deviceID)
	if err != nil {
		log.Printf("Device deactivation failed for %s: %v", deviceID, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("Device deactivation failed: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func getDeviceConfig(c *gin.Context) {
	config := setting.GetThingsConfig()
	log.Printf("Fetching device config: %+v", config)
	c.JSON(http.StatusOK, SuccessResponse(config))
}

func generateGetStorageCredentials(c *gin.Context) {
	deviceId := getDeviceID(c)
	if deviceId == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse("Missing device-id header"))
		return
	}

	// Get storage credentials
	actions := []external_storage.Action{external_storage.GetObject, external_storage.PutObject}
	credentials, err := cloudStorageService.GetStorageCredentials(actions, deviceId)
	if err != nil {
		log.Printf("Failed to get storage credentials for device %s: %v", deviceId, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("Failed to get storage credentials: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(credentials))
}

func signIn(c *gin.Context) {
	var req user.SignInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request format: "+err.Error()))
		return
	}

	signInResponse, err := userService.SignIn(req)
	if err != nil {
		log.Printf("Sign-in failed for account %s: %v", req.Account, err)
		c.JSON(http.StatusUnauthorized, ErrorResponse("Sign-in failed: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(signInResponse))
}

func signUp(c *gin.Context) {
	var req user.SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request format: "+err.Error()))
		return
	}

	signUpResponse, err := userService.SignUp(req)
	if err != nil {
		log.Printf("Sign-up failed for account %s: %v", req.Account, err)
		c.JSON(http.StatusConflict, ErrorResponse("Sign-up failed: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse(signUpResponse))
}

func userDetails(c *gin.Context) {
	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse("Missing user-id header"))
		return
	}

	log.Printf("Fetching details for user ID: %s", userID)

	// Get user details
	user, err := userService.Get(userID)
	if err != nil {
		log.Printf("Error fetching user %s: %v", userID, err)
		c.JSON(http.StatusNotFound, ErrorResponse("User not found"))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(user))
}

func handleEvent(c *gin.Context) {
	deviceID := getDeviceID(c)
	var event things.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request payload"))
		return
	}

	// Process the device event
	if err := thingsService.HandleEvent(deviceID, event); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse("Failed to report device event"))
		return
	}

	c.JSON(http.StatusOK, nil)
}

func handleState(c *gin.Context) {
	deviceID := getDeviceID(c)
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse("Missing device-id header"))
		return
	}

	var state things.State
	if err := c.ShouldBindJSON(&state); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request payload"))
		return
	}

	// Process the device state
	if err := thingsService.HandleState(deviceID, state); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse("Failed to report device state"))
		return
	}
	c.JSON(http.StatusOK, nil)
}

func fetchDesiredState(c *gin.Context) {
	deviceID := getDeviceID(c)
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse("Missing device-id header"))
		return
	}

	// Get device desired state
	state, err := thingsService.FetchDesiredState(deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse("Failed to get device desired state"))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(state))
}
