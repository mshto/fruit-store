package database

import (
	"testing"

	"github.com/kr/pgtest"
	"github.com/stretchr/testify/assert"
)

func TestNewDatabasePositive(t *testing.T) {
	pg := pgtest.Start(t)
	defer pg.Stop()

	_, err := New(Database{
		DBType: "postgres",
	})
	assert.NotNil(t, err)
}

func TestNewDatabaseNagative(t *testing.T) {
	_, err := New(Database{})
	assert.NotNil(t, err)
}
