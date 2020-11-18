package logger

import (
	"github.com/sirupsen/logrus"
	logrusPrefixed "github.com/x-cray/logrus-prefixed-formatter"
)

//go:generate mockgen -destination=./mocks/logger.go -package=mocks -source=./logger.go Logger

//Logger struct stores system state logger configuration
type Logger struct {
	LogLevel string `json:"LogLevel" envconfig:"LogLevel" validate:"required"`
}

//New is a function for opening and loading log file
func New(logger Logger) (*logrus.Logger, error) {
	logrusLogLevel, err := logrus.ParseLevel(logger.LogLevel)
	if err != nil {
		return nil, err
	}

	log := logrus.New()
	log.Level = logrusLogLevel
	log.Formatter = &logrusPrefixed.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	}

	return log, err
}
