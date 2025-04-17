package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/config"
	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/arefev/gophkeeper/internal/client/tui/step"
	"github.com/arefev/gophkeeper/internal/logger"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conf, err := config.NewConfig(os.Args[1:])
	if err != nil {
		return fmt.Errorf("run: init config failed: %w", err)
	}

	l, err := logger.Build(conf.LogLevel)
	if err != nil {
		return fmt.Errorf("run: logger build failed: %w", err)
	}

	conn := connection.NewGRPCClient(l)
	if err = conn.Connect(conf.Address); err != nil {
		return fmt.Errorf("run: connect to server failed: %w", err)
	}

	err = conn.FileUpload(context.Background(), []byte("new file data"))
	if err != nil {
		l.Error("FileUpload failed", zap.Error(err))
	}

	defer func() {
		if err = conn.Close(); err != nil {
			l.Error("connect close failed", zap.Error(err))
		}
	}()

	if conf.LogFilePath != "" {
		f, err := tea.LogToFile(conf.LogFilePath, conf.LogLevel)
		if err != nil {
			return fmt.Errorf("run: log to file failed: %w", err)
		}
		defer func() {
			if err = f.Close(); err != nil {
				l.Error("log file close failed: %w", zap.Error(err))
			}
		}()
	}

	_, err = step.NewStart(app.NewApp(conn, l)).NewProgram().Run()
	if err != nil {
		return fmt.Errorf("run: app stopped with error: %w", err)
	}

	return nil
}
