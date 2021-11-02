package config

import (
	"errors"
	"log"
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

//Config holds configuration for the project
type Config struct {
	Port           string `env:"PORT, default=6666"`
	Env            string `env:"ENV, default=development"`
	Database       DatabaseConfig
	JWTConfig      JWTConfig
	InternalConfig InternalConfig
}

//DatabaseConfig holds configuration for database
type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST, default=localhost"`
	Port     string `env:"DATABASE_PORT, default=5432"`
	Username string `env:"DATABASE_USERNAME, required"`
	Password string `env:"DATABASE_PASSWORD, required"`
	Name     string `env:"DATABASE_Name, required"`
}

//JWTConfig holds configuration for JWT secret key
type JWTConfig struct {
	SecretKey string `env:"JWT_SECRET_KEY, required"`
}

//InternalConfig holds configuration internal communication between microservices
type InternalConfig struct {
	Username string `env:"SVC_USERNAME, required"`
	Password string `env:"SVC_PASSWORD, required"`
}

func NewConfig(env string) (*Config, error) {
	var config Config
	if err := godotenv.Load(env); err != nil {
		log.Println(errors.Unwrap(err))
	}

	if err := envdecode.Decode(&config); err != nil {
		return nil, errors.Unwrap(err)
	}

	return &config, nil
}