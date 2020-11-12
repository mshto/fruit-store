package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/go-playground/validator.v9"

	"github.com/mshto/fruit-store/cache"
	"github.com/mshto/fruit-store/database"
	"github.com/mshto/fruit-store/logger"
)

//Config struct stores system state configuration
type Config struct {
	ListenURL  string            `json:"ListenURL"     envconfig:"PORT"     validate:"required"`
	URLPrefix  string            `json:"URLPrefix"     envconfig:"URLPrefix"     validate:"required"`
	APIVersion string            `json:"APIVersion"`
	Logger     logger.Logger     `json:"Logger"`
	Database   database.Database `json:"Database"`
	Redis      cache.Redis       `json:"Redis"`
	Auth       Auth              `json:"Auth"`
	Sales      []GeneralSale
}

// Auth struct stores auth secret keys
type Auth struct {
	AccessSecret  string `json:"AccessSecret"    envconfig:"AUTH_ACCESS_SECRET"     validate:"required"`
	RefreshSecret string `json:"RefreshSecret"   envconfig:"AUTH_REFRESH_SECRET"    validate:"required"`
}

// GeneralSale GeneralSale
type GeneralSale struct {
	ID       string
	Elements map[string]int
	Rule     string
	Discount int
}

// New is reading json file, validating and returning config
func New(configPath string, salesCfg string) (*Config, error) {
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

	salesContents, err := ioutil.ReadFile(salesCfg)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(salesContents, &config.Sales)
	if err != nil {
		return nil, err
	}

	readConfigFromENV(config)
	// getEnvVariables(config)
	return config, nil
}

// validate validates a struct and nested fields
func validate(c *Config) error {
	v := validator.New()

	return v.Struct(c)
}

// readConfigFromENV reads data from environment variables
func readConfigFromENV(cfg *Config) error {
	return envconfig.Process("", cfg)

}

// func getEnvVariables(c *Config) {
// 	port := os.Getenv("PORT")
// 	if port != "" {
// 		c.ListenURL = ":" + port
// 	}
// 	as := os.Getenv("ACCESS_SECRET")
// 	if as != "" {
// 		c.Auth.AccessSecret = as
// 	}
// 	rs := os.Getenv("REFRESH_SECRET")
// 	if rs != "" {
// 		c.Auth.RefreshSecret = rs
// 	}
// }
