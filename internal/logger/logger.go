package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Build(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("zap logger parse level failed: %w", err)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("zap logger build from config failed: %w", err)
	}

	return zl, nil
}

func InFile(path, level string) (*zap.Logger, error) {
	writeSyncer, err := os.Create(path) // Log file storage directory
	if err != nil {
		return nil, fmt.Errorf("zap logger create file failed: %w", err)
	}
	encoderConfig := zap.NewProductionEncoderConfig() // Specify time format
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)               // Get the encoder, NewJSONEncoder() outputs in JSON format, NewConsoleEncoder() outputs in plain text format
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel) // The third and subsequent parameters are the log levels for writing to the file. In ErrorLevel mode, only error - level logs are recorded.
	return zap.New(core, zap.AddCaller()), nil
}
