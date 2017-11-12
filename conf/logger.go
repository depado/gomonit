package conf

import (
	"github.com/sirupsen/logrus"
)

// Logger is a configuration struct to define logger's behaviour
type Logger struct {
	Level  string `yaml:"level" default:"info"`
	Format string `yaml:"format" default:"text"`
}

// Configure takes the configuration for the logger and translats it to
// logrus's usage
func (c Logger) Configure() {
	SetFormatter(c.Format)
	SetLogLevel(c.Level)
}

// SetLogLevel sets the logging level when possible, otherwise it fallbacks to
// the default logrus level and logs a warning
func SetLogLevel(lvl string) {
	l, err := logrus.ParseLevel(lvl)
	if err != nil {
		logrus.WithField("provided", lvl).Warn("Invalid log level, fallback to Info level")
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(l)
	}
}

// SetFormatter defines the way logs are formatted
func SetFormatter(format string) {
	switch format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
}
