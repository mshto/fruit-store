package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultLoggerPossitive(t *testing.T) {
	conf := Logger{
		LogLevel: "info",
	}
	log, err := New(conf)

	assert.NotNil(t, log)
	assert.Nil(t, err)
}

func TestWrongLogLevelNegative(t *testing.T) {
	conf := Logger{
		LogLevel: "wrong level",
	}
	_, err := New(conf)

	assert.NotNil(t, err)
}
