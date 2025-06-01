// Package env provides utilities for working with environment variables.
// It includes functions for retrieving and parsing environment variables
// with type safety and fallback values.
package env

import (
	"os"
	"strconv"
	"time"
)

// GetString retrieves a string value from the environment.
// If the environment variable is not set, it returns the fallback value.
//
// Example:
//
//	port := env.GetString("PORT", "8080")
func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}
	return val
}

// GetInt retrieves an integer value from the environment.
// If the environment variable is not set or cannot be parsed as an integer,
// it returns the fallback value.
//
// Example:
//
//	maxConns := env.GetInt("MAX_CONNECTIONS", 100)
func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return num
}

// GetDuration retrieves a time.Duration value from the environment.
// If the environment variable is not set or cannot be parsed as a duration,
// it returns the fallback value.
//
// Example:
//
//	timeout := env.GetDuration("TIMEOUT", 5*time.Second)
func GetDuration(key string, fallback time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	duration, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}

	return duration
}
