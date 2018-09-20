package main

import (
	"log"
	"os"
)

// GetEnv : the .env variable by key or its default
func GetEnv(key, fallback string) string {
	log.Printf("looking for env key: %s", key)
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return fallback
}
