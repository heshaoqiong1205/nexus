package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// Point represents a geographic coordinate point for PostgreSQL POINT type
type Point struct {
	Latitude  float64 `json:"latitude"`  // Latitude
	Longitude float64 `json:"longitude"` // Longitude
}

// String returns the string representation of the point
func (p Point) String() string {
	// PostgreSQL POINT format: (longitude, latitude)
	return fmt.Sprintf("(%f,%f)", p.Longitude, p.Latitude)
}

// Value implements the driver.Valuer interface for database storage
func (p Point) Value() (driver.Value, error) {
	return p.String(), nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (p *Point) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("cannot scan %T into Point", value)
	}

	// Parse PostgreSQL point format: (x,y)
	str = strings.Trim(str, "()")
	parts := strings.Split(str, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid point format: %s", str)
	}

	x, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return fmt.Errorf("invalid x coordinate: %s", parts[0])
	}

	y, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return fmt.Errorf("invalid y coordinate: %s", parts[1])
	}

	// In PostgreSQL POINT type: x = longitude, y = latitude
	p.Longitude = x
	p.Latitude = y
	return nil
}

// NewPoint creates a new Point with the given coordinates
func NewPoint(latitude, longitude float64) *Point {
	return &Point{Latitude: latitude, Longitude: longitude}
}
