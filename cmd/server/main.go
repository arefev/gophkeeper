package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/config"
	"github.com/arefev/gophkeeper/internal/server/service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	conf, err := config.NewConfig(os.Args[1:])
	if err != nil {
		return fmt.Errorf("run: init config fail: %w", err)
	}

	databaseDSN := databaseDSN(&conf)
	err = migrationsUp(databaseDSN)
	if err != nil {
		return fmt.Errorf("run: migration up fail: %w", err)
	}

	err = runGRPC(ctx, &conf)
	if err != nil {
		return fmt.Errorf("run: runGRPC fail: %w", err)
	}

	return nil
}

func runGRPC(ctx context.Context, c *config.Config) error {
	listen, err := net.Listen("tcp", c.Address)
	if err != nil {
		return fmt.Errorf("runGRPC Listen failed: %w", err)
	}

	s := grpc.NewServer()
	proto.RegisterRegistrationServer(s, &service.GRPCServer{})

	go func() {
		<-ctx.Done()
		fmt.Println("Server stopped")
		s.Stop()
	}()

	fmt.Println("Server running")

	if err := s.Serve(listen); err != nil {
		return fmt.Errorf("runGRPC Serve failed: %w", err)
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

func databaseDSN(cnf *config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cnf.DBName,
		cnf.DBPassword,
		cnf.DBHost,
		cnf.DBPort,
		cnf.DBName,
	)
}
