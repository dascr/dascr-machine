package config

import (
	"os"

	"github.com/dascr/dascr-machine/service/logger"
)

// MustGet will check for the env settings
func MustGet(k string) string {
	v := os.Getenv(k)
	if v == "" {
		logger.MissingEnv(k)
	}
	return v
}
