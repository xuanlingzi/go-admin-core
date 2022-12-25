package config

import "github.com/xuanlingzi/go-admin-core/sdk/pkg/logger"

type Logger struct {
	Type      string `json:"type,omitempty" yaml:"type"`
	Path      string `json:"path,omitempty" yaml:"path"`
	Level     string `json:"level,omitempty" yaml:"level"`
	Stdout    string `json:"stdout,omitempty" yaml:"stdout"`
	EnabledDB bool   `json:"enabled_db,omitempty" yaml:"enabled_db"`
	Cap       uint   `json:"cap,omitempty" yaml:"cap"`
}

// Setup 设置logger
func (e Logger) Setup() {
	logger.SetupLogger(
		logger.WithType(e.Type),
		logger.WithPath(e.Path),
		logger.WithLevel(e.Level),
		logger.WithStdout(e.Stdout),
		logger.WithCap(e.Cap),
	)
}

var LoggerConfig = new(Logger)
