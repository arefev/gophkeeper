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
	"github.com/arefev/gophkeeper/internal/server/handler"
	"github.com/arefev/gophkeeper/internal/server/handler/interceptor"
	"github.com/arefev/gophkeeper/internal/server/repository"
	"github.com/arefev/gophkeeper/internal/server/trm"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	log.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
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

	intr := interceptor.New(app)
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			intr.UnaryCheckToken(),
		),
		grpc.ChainStreamInterceptor(
			intr.StreamCheckToken(),
		),
	)
	proto.RegisterAuthServer(s, handler.NewAuthHandler(app))
	proto.RegisterFileServer(s, handler.NewFileHandler(app))
	proto.RegisterListServer(s, handler.NewListHandler(app))

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
	host := net.JoinHostPort(cnf.DBHost, cnf.DBPort)
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		cnf.DBName,
		cnf.DBPassword,
		host,
		cnf.DBName,
	)
}
