package main

import (
	"fmt"

	"nexus/pkg/setting"
)

func main() {
	fmt.Println("Testing Nexus Configuration Loading...")
	fmt.Println("=====================================")

	// Load configuration
	setting.Setup()

	// Display configuration
	fmt.Printf("App Section:\n")
	fmt.Printf("  Name: %s\n", setting.AppSetting.Name)

	fmt.Printf("\nServer Section:\n")
	fmt.Printf("  RunMode: %s\n", setting.ServerSetting.RunMode)
	fmt.Printf("  HttpPort: %d\n", setting.ServerSetting.HttpPort)
	fmt.Printf("  ReadTimeout: %s\n", setting.ServerSetting.ReadTimeout)
	fmt.Printf("  WriteTimeout: %s\n", setting.ServerSetting.WriteTimeout)

	fmt.Printf("\nDatabase Section:\n")
	fmt.Printf("  Type: %s\n", setting.DatabaseSetting.Type)
	fmt.Printf("  Host: %s\n", setting.DatabaseSetting.Host)
	fmt.Printf("  Port: %d\n", setting.DatabaseSetting.Port)
	fmt.Printf("  User: %s\n", setting.DatabaseSetting.User)
	fmt.Printf("  Password: %s\n", setting.DatabaseSetting.Password)
	fmt.Printf("  Name: %s\n", setting.DatabaseSetting.Name)
	fmt.Printf("  TablePrefix: %s\n", setting.DatabaseSetting.TablePrefix)

	fmt.Printf("\nRedis Section:\n")
	fmt.Printf("  Host: %s\n", setting.RedisSetting.Host)
	fmt.Printf("  Password: %s\n", setting.RedisSetting.Password)
	fmt.Printf("  MaxIdle: %d\n", setting.RedisSetting.MaxIdle)
	fmt.Printf("  MaxActive: %d\n", setting.RedisSetting.MaxActive)
	fmt.Printf("  IdleTimeout: %s\n", setting.RedisSetting.IdleTimeout)

	fmt.Println("\n✓ Configuration loaded successfully!")
}
