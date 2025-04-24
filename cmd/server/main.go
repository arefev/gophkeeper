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
	"github.com/arefev/gophkeeper/internal/server/handler/interceptor"

	// "github.com/arefev/gophkeeper/internal/server/service"

	// "github.com/arefev/gophkeeper/internal/server/model"
	"github.com/arefev/gophkeeper/internal/server/handler"
	"github.com/arefev/gophkeeper/internal/server/repository"
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

	// err = app.TrManager.Do(ctx, func(ctx context.Context) error {
	// 	meta, err := app.Rep.Meta.Find(ctx, 28)
	// 	if err != nil {
	// 		return fmt.Errorf("run: meta get failed: %w", err)
	// 	}
	// 	// l.Sugar().Infof("meta %+v", meta)

	// 	es := service.NewEncryptionService(app)
	// 	data, err := es.Decrypt(meta.File.Data)
	// 	if err != nil {
	// 		return fmt.Errorf("run: decrypt data failed: %w", err)
	// 	}

	// 	// l.Sugar().Infof("data %+v", string(data))
	// 	file, err := os.Create("./" + meta.File.Name)
	// 	if err != nil {
	// 		return fmt.Errorf("run: create file failed: %w", err)
	// 	}

	// 	_, err = file.Write(data)
	// 	if err != nil {
	// 		return fmt.Errorf("run: write file failed: %w", err)
	// 	}
	// 	file.Close()

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

	intr := interceptor.New(app)
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(intr.UnaryCheckToken()),
		),
		grpc.ChainStreamInterceptor(
			grpc.StreamServerInterceptor(intr.StreamCheckToken()),
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
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cnf.DBName,
		cnf.DBPassword,
		cnf.DBHost,
		cnf.DBPort,
		cnf.DBName,
	)
}
