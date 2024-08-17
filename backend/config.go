package main

import (
	"fmt"
	"os"
)

type Config struct {
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
	isLocal    bool
}

func NewConfigFromEnv() (*Config, error) {
	dbHost, exists := os.LookupEnv("ALBATROSS_DB_HOST")
	if !exists {
		return nil, fmt.Errorf("ALBATROSS_DB_HOST not set")
	}
	dbPort, exists := os.LookupEnv("ALBATROSS_DB_PORT")
	if !exists {
		return nil, fmt.Errorf("ALBATROSS_DB_PORT not set")
	}
	dbUser, exists := os.LookupEnv("ALBATROSS_DB_USER")
	if !exists {
		return nil, fmt.Errorf("ALBATROSS_DB_USER not set")
	}
	dbPassword, exists := os.LookupEnv("ALBATROSS_DB_PASSWORD")
	if !exists {
		return nil, fmt.Errorf("ALBATROSS_DB_PASSWORD not set")
	}
	dbName, exists := os.LookupEnv("ALBATROSS_DB_NAME")
	if !exists {
		return nil, fmt.Errorf("ALBATROSS_DB_NAME not set")
	}
	isLocalStr, exists := os.LookupEnv("ALBATROSS_IS_LOCAL")
	isLocal := exists && isLocalStr == "1"
	return &Config{
		dbHost:     dbHost,
		dbPort:     dbPort,
		dbUser:     dbUser,
		dbPassword: dbPassword,
		dbName:     dbName,
		isLocal:    isLocal,
	}, nil
}
