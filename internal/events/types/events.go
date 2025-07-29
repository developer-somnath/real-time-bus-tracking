package types

import "time"

type BusCreatedEvent struct {
	BusID             int       `json:"bus_id"`
	LicensePlate      string    `json:"license_plate"`
	WheelchairEnabled bool      `json:"wheelchair_enabled"`
	CreatedAt         time.Time `json:"created_at"`
}

type DriverLocationUpdatedEvent struct {
	DriverID  int       `json:"driver_id"`
	BusID     int       `json:"bus_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}
