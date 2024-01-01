package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/benharmonics/backend/logging"
)

const (
	AppEnvironment = "APP_ENV"
	Development    = "development"
)

func ValidateConfig() error {
	var errors []error
	errors = append(errors, NewMongoConfig().Validate()...)
	errors = append(errors, NewAppConfig().Validate()...)
	if len(errors) == 0 {
		return nil
	}
	for _, err := range errors {
		logging.Error(err)
	}
	return fmt.Errorf("misconfigured environment variables")
}

func getEnv(key, defaultval string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultval
}

func getEnvAsInt(key string, defaultval int) int {
	if val := os.Getenv(key); val != "" {
		if iVal, err := strconv.Atoi(val); err == nil {
			return iVal
		}
	}
	return defaultval
}
