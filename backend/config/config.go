package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gopkg.in/go-playground/validator.v9"

	"github.com/mshto/fruit-store/backend/cache"
	"github.com/mshto/fruit-store/backend/database"
	"github.com/mshto/fruit-store/backend/logger"
)

//Config struct stores system state configuration
type Config struct {
	ListenURL  string `validate:"required"`
	URLPrefix  string `validate:"required"`
	APIVersion string
	Logger     logger.Logger
	Database   database.Database
	Redis      cache.Redis
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
	return config, nil
}

// validate validates a struct and nested fields
func validate(c *Config) error {
	v := validator.New()

	return v.Struct(c)
}
