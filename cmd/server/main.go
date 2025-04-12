package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/arefev/gophkeeper/internal/server/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conf, err := config.NewConfig(os.Args[1:])
	if err != nil {
		return fmt.Errorf("run: init config fail: %w", err)
	}

	err = migrationsUp(conf.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("run: migration up fail: %w", err)
	}

	return nil
}

func migrationsUp(dsn string) error {
	ex, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path fail: %w", err)
	}

	filePath := filepath.Dir(ex)
	m, err := migrate.New("file://"+filePath+"/db/migrations", dsn)
	if err != nil {
		return fmt.Errorf("migrations instance fail: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrations up fail: %w", err)
	}

	return nil
}
