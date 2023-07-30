package api

import "time"

type BalenaServerInfo struct {
	BaseURL  string
	APIToken string
}

type Device struct {
	UUID       string `json:"uuid" bson:"uuid"`
	ID         int    `json:"id" bson:"id"`
	IsOnline   bool   `json:"is_online" bson:"is_online"`
	Status     string `json:"status" bson:"status"`
	DeviceName string `json:"device_name" bson:"device_name"`

	Longitude *string `json:"longitude" bson:"longitude"`
	Latitude  *string `json:"latitude" bson:"latitude"`
	Location  *string `json:"location" bson:"location"`

	LastConnectivityEvent *time.Time `json:"last_connectivity_event" bson:"last_connectivity_event"`
	MemoryUsage           int        `json:"memory_usage" bson:"memory_usage"`
	MemoryTotal           int        `json:"memory_total" bson:"memory_total"`
	StorageUsage          int        `json:"storage_usage" bson:"storage_usage"`
	StorageTotal          int        `json:"storage_total" bson:"storage_total"`
	CPUUsage              int        `json:"cpu_usage" bson:"cpu_usage"`
	CPUTemp               int        `json:"cpu_total" bson:"cpu_total"`
	IsUndervolted         bool       `json:"is_undervolted" bson:"is_undervolted"`
}
