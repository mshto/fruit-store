package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/go-playground/validator.v9"

	"github.com/mshto/fruit-store/cache"
	"github.com/mshto/fruit-store/database"
	"github.com/mshto/fruit-store/logger"
)

//Config struct stores system state configuration
type Config struct {
	ListenURL  string `validate:"required"`
	URLPrefix  string `validate:"required"`
	APIVersion string
	Logger     logger.Logger
	Database   database.Database
	Redis      cache.Redis
	Auth       Auth
}

// Auth struct stores auth secret keys
type Auth struct {
	AccessSecret  string
	RefreshSecret string
}

// New is reading json file, validating and returning config
func New(configPath string) (*Config, error) {
	config := new(Config)
	contents, err := ioutil.ReadFile(configPath) // nolint: gosec
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(contents, config)
	if err != nil {
		return nil, err
	}
	err = validate(config)
	if err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	getEnvVariables(config)
	return config, nil
}

// validate validates a struct and nested fields
func validate(c *Config) error {
	v := validator.New()

	return v.Struct(c)
}

func getEnvVariables(c *Config) {
	port := os.Getenv("PORT")
	if port != "" {
		c.ListenURL = ":" + port
	}
	as := os.Getenv("ACCESS_SECRET")
	if as != "" {
		c.Auth.AccessSecret = as
	}
	rs := os.Getenv("REFRESH_SECRET")
	if rs != "" {
		c.Auth.RefreshSecret = rs
	}
}
