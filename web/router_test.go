package web

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	redismock "github.com/mshto/fruit-store/cache/mock"
	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/repository"
	loggermock "github.com/sirupsen/logrus/hooks/test"
)

func TestConfigValidation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger, _ := loggermock.NewNullLogger()

	route := New(&config.Config{}, logger, repository.New(db), redismock.NewMockCache(mockCtrl))
	assert.NotNil(t, route)
}
