package utils

import (
	"os"
)

// GetEnv retrieves an environment variable value with a fallback default
// if the variable is not set.
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
