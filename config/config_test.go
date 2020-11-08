package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidation(t *testing.T) {
	empty := &Config{}

	err := validate(empty)
	assert.NotNil(t, err)
}

func TestNewEmptyFile(t *testing.T) {
	_, err := New("")
	assert.Error(t, err, "expected error on empty path")
}

func TestNewDefaultConfig(t *testing.T) {
	_, err := New("mock/valid_config.json")
	assert.Nil(t, err)
}

func TestNewInvalidJsonConfig(t *testing.T) {
	_, err := New("mock/invalid_config.json")
	assert.NotNil(t, err)
}

func TestNewInvalidJsonValidationConfig(t *testing.T) {
	_, err := New("mock/validation_err_config.json")
	assert.NotNil(t, err)
}
