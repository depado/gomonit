package conf

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSetLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		lvl      string
		expected logrus.Level
	}{
		{"must set to debug", "debug", logrus.DebugLevel},
		{"must set to info", "info", logrus.InfoLevel},
		{"must set to warning", "warn", logrus.WarnLevel},
		{"must set to error", "error", logrus.ErrorLevel},
		{"must not fail and fallback", "random", logrus.InfoLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLogLevel(tt.lvl)
			assert.Equal(t, logrus.GetLevel(), tt.expected)
		})
	}
}

func TestSetFormatter(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name string
		args args
		want logrus.Formatter
	}{
		{"text format", args{"text"}, &logrus.TextFormatter{}},
		{"json format", args{"json"}, &logrus.JSONFormatter{}},
		{"invalid format", args{"invalid"}, &logrus.TextFormatter{}},
		{"empty arg", args{""}, &logrus.TextFormatter{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetFormatter(tt.args.format)
			assert.Equal(t, tt.want, logrus.StandardLogger().Formatter)
		})
	}
}

func TestLogger_Configure(t *testing.T) {
	type fields struct {
		Level  string
		Format string
	}
	type want struct {
		Formatter logrus.Formatter
		Level     logrus.Level
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{"empty", fields{}, want{&logrus.TextFormatter{}, logrus.InfoLevel}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Logger{
				Level:  tt.fields.Level,
				Format: tt.fields.Format,
			}
			c.Configure()
			assert.Equal(t, tt.want.Level, logrus.GetLevel())
		})
	}
}
