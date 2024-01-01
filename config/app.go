package config

import "github.com/benharmonics/backend/utils/validation"

type AppConfig struct {
	Host string `validate:"required"`
	Port int    `validate:"required,gte=1024,lte=49151"`
}

func NewAppConfig() AppConfig {
	return AppConfig{
		Host: getEnv("HOST", "localhost"),
		Port: getEnvAsInt("PORT", 9090),
	}
}

func (conf AppConfig) Validate() []error {
	return validation.ValidateStructAll(conf, nil)
}
