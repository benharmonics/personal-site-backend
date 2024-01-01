package config

import "github.com/benharmonics/backend/utils/validation"

type MongoConfig struct {
	Username string `validate:"required_with=Password"`
	Password string `validate:"required_with=Username"`
	Proto    string `validate:"required,oneof=mongodb mongodb+srv"`
	Host     string `validate:"required"`
	Port     int    `validate:"required,gte=1024,lte=49151"`
}

func NewMongoConfig() MongoConfig {
	return MongoConfig{
		Username: getEnv("MONGODB_USERNAME", ""),
		Password: getEnv("MONGODB_PASSWORD", ""),
		Proto:    getEnv("MONGODB_PROTO", "mongodb"),
		Host:     getEnv("MONGODB_HOST", "localhost"),
		Port:     getEnvAsInt("MONGODB_PORT", 27017),
	}
}

func (conf MongoConfig) Validate() []error {
	return validation.ValidateStructAll(conf, nil)
}
