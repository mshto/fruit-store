package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	logrusPrefixed "github.com/x-cray/logrus-prefixed-formatter"
)

//Logger struct stores system state logger configuration
type Logger struct {
	Path     string `validate:"required"`
	LogLevel string `validate:"required"`
}

//New is a function for opening and loading log file
func New(logger Logger) (*logrus.Logger, error) {
	file, err := os.OpenFile(logger.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	logrusLogLevel, err := logrus.ParseLevel(logger.LogLevel)
	if err != nil {
		return nil, err
	}

	return &logrus.Logger{
		Out:   file,
		Level: logrusLogLevel,
		Formatter: &logrusPrefixed.TextFormatter{
			DisableColors:   false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}, err
}
