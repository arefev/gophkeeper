package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/config"
	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/arefev/gophkeeper/internal/client/tui/step"
	"github.com/arefev/gophkeeper/internal/logger"
	"go.uber.org/zap"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	conf, err := config.NewConfig(os.Args[1:])
	if err != nil {
		return fmt.Errorf("run: init config failed: %w", err)
	}

	l, err := createLogger(conf)
	if err != nil {
		return fmt.Errorf("run: logger failed: %w", err)
	}

	conn := connection.NewGRPCClient(conf.ChunkSize, l)
	if err = conn.Connect(conf.Address); err != nil {
		return fmt.Errorf("run: connect to server failed: %w", err)
	}

	defer func() {
		if err = conn.Close(); err != nil {
			l.Error("connect close failed", zap.Error(err))
		}
	}()

	_, err = step.NewStart(app.NewApp(conn, l)).NewProgram().Run()
	if err != nil {
		return fmt.Errorf("run: app stopped with error: %w", err)
	}

	return nil
}

func createLogger(conf *config.Config) (*zap.Logger, error) {
	var l *zap.Logger
	var err error
	if conf.LogFilePath != "" {
		l, err = logger.InFile(conf.LogFilePath, conf.LogLevel)
		if err != nil {
			return nil, fmt.Errorf("create logger in file failed: %w", err)
		}
	} else {
		l, err = logger.Build(conf.LogLevel)
		if err != nil {
			return nil, fmt.Errorf("create logger build failed: %w", err)
		}
	}

	return l, nil
}
