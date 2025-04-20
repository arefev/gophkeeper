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

	"github.com/arefev/gophkeeper/internal/logger"
	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/application"
	"github.com/arefev/gophkeeper/internal/server/config"
	"github.com/arefev/gophkeeper/internal/server/db/postgresql"
	// "github.com/arefev/gophkeeper/internal/server/model"
	"github.com/arefev/gophkeeper/internal/server/repository"
	"github.com/arefev/gophkeeper/internal/server/server"
	"github.com/arefev/gophkeeper/internal/server/trm"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
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

	l, err := logger.Build(conf.LogLevel)
	if err != nil {
		return fmt.Errorf("run: logger build fail: %w", err)
	}

	databaseDSN := databaseDSN(conf)
	db, err := postgresql.NewDB(l).Connect(databaseDSN)
	if err != nil {
		return fmt.Errorf("run: db trm connect fail: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			l.Error("db close failed: %w", zap.Error(err))
		}
	}()

	err = migrationsUp(databaseDSN)
	if err != nil {
		return fmt.Errorf("run: migration up fail: %w", err)
	}

	tr := trm.NewTr(db.Connection())
	app := &application.App{
		Rep: application.Repository{
			User: repository.NewUser(tr, l),
			Meta: repository.NewMeta(tr, l),
		},
		TrManager: trm.NewTrm(tr, l),
		Log:       l,
		Conf:      conf,
	}

	// path := "/home/arefev/dev/study/golang/gophkeeper/storage/405524d0-53a7-4a6d-8eee-4914b06b163e/office.iso"
	// path := "/home/arefev/dev/study/golang/gophkeeper/storage/5372edba-7376-4cc1-88bc-153d2b93421c/test.txt"
	// path := "/home/arefev/dev/study/golang/gophkeeper/storage/2dc22093-61b4-40fd-9dea-9dacd5cff4c1/joxi.exe"
	// file, err := os.ReadFile(path)
	// if err != nil {
	// 	return fmt.Errorf("run: read file failed: %w", err)
	// }

	// err = app.TrManager.Do(ctx, func(ctx context.Context) error {
	// 	err = app.Rep.Meta.Create(ctx, &model.Meta{
	// 		UserID: 1,
	// 		Type: model.MetaTypeFile,
	// 		Name: "test",
	// 	}, &model.File{Name: "test.txt", Data: file})
	// 	if err != nil {
	// 		return fmt.Errorf("run: meta create failed: %w", err)
	// 	}

	// 	return nil
	// })
	// if err != nil {
	// 	return fmt.Errorf("run: do transaction failed: %w", err)
	// }

	err = runServer(ctx, app, conf, l)
	if err != nil {
		return fmt.Errorf("run: runGRPC fail: %w", err)
	}

	return nil
}

func runServer(ctx context.Context, app *application.App, c *config.Config, l *zap.Logger) error {
	listen, err := net.Listen("tcp", c.Address)
	if err != nil {
		return fmt.Errorf("runGRPC Listen failed: %w", err)
	}

	s := grpc.NewServer()
	proto.RegisterAuthServer(s, server.NewAuthServer(app))
	proto.RegisterFileServer(s, server.NewFileServer(app))

	go func() {
		<-ctx.Done()
		l.Info("Server stopped")
		s.Stop()
	}()

	l.Info("Server running", zap.String("port", c.Address), zap.String("log level", c.LogLevel))

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
