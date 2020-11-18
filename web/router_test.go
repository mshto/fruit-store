package web

import (
	"testing"

	"github.com/mshto/fruit-store/cache"
	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/repository"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestConfigValidation(t *testing.T) {
	route := New(&config.Config{}, &logrus.Logger{}, &repository.Repository{}, &cache.Cache{})
	assert.NotNil(t, route)
}
