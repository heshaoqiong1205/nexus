package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"nexus/pkg/setting"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting Nexus Server with Configuration...")
	fmt.Println("==========================================")

	// Load configuration
	setting.Setup()

	// Set Gin mode
	gin.SetMode(setting.ServerSetting.RunMode)

	// Create router
	router := gin.Default()

	// Add basic routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Nexus Things Platform API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	router.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"app": gin.H{
				"name": setting.AppSetting.Name,
			},
			"server": gin.H{
				"runMode":      setting.ServerSetting.RunMode,
				"httpPort":     setting.ServerSetting.HttpPort,
				"readTimeout":  setting.ServerSetting.ReadTimeout.String(),
				"writeTimeout": setting.ServerSetting.WriteTimeout.String(),
			},
			"database": gin.H{
				"type":        setting.DatabaseSetting.Type,
				"host":        setting.DatabaseSetting.Host,
				"port":        setting.DatabaseSetting.Port,
				"name":        setting.DatabaseSetting.Name,
				"tablePrefix": setting.DatabaseSetting.TablePrefix,
			},
		})
	})

	// Configure server
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("✓ Configuration loaded successfully!\n")
	fmt.Printf("✓ Server starting on port %d\n", setting.ServerSetting.HttpPort)
	fmt.Printf("✓ Database: %s@%s:%d/%s\n",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Port,
		setting.DatabaseSetting.Name)
	fmt.Printf("✓ Visit: http://localhost:%d\n\n", setting.ServerSetting.HttpPort)

	log.Printf("Starting server on port %d", setting.ServerSetting.HttpPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
