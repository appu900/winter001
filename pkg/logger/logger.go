package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(serviceName string, env string, level string) (*zap.Logger, error) {
	var cfg zap.Config
	if env == "development" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "ts"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	if level != "" {
		lvl, err := zapcore.ParseLevel(level)
		if err != nil {
			return nil, err
		}
		cfg.Level = zap.NewAtomicLevelAt(lvl)
	}

	cfg.InitialFields = map[string]any{
		"service": serviceName,
		"env":     env,
	}

	return cfg.Build(zap.AddCaller())
}
