package web

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	redismock "github.com/mshto/fruit-store/cache/mocks"
	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/repository"
)

func TestConfigValidation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	route := New(&config.Config{}, &logrus.Logger{}, &repository.Repository{}, redismock.NewMockCache(mockCtrl))
	assert.NotNil(t, route)
}
