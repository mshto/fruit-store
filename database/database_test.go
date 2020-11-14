package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDatabasePositive(t *testing.T) {
	_, err := New(Database{
		DBType: "postgres",
	})
	assert.NotNil(t, err)
}

func TestNewDatabaseNagative(t *testing.T) {
	_, err := New(Database{})
	assert.NotNil(t, err)
}
