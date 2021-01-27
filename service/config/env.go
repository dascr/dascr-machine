package config

import (
	"log"
	"os"
)

// MustGet will check for the env settings
func MustGet(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("Missing environment variable '%+v'", k)
	}
	return v
}
