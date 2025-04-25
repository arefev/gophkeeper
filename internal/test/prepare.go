package test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/arefev/gophkeeper/internal/logger"
	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/application"
	"github.com/arefev/gophkeeper/internal/server/config"
	"github.com/arefev/gophkeeper/internal/server/db/postgresql"
	"github.com/arefev/gophkeeper/internal/server/handler"
	"github.com/arefev/gophkeeper/internal/server/handler/interceptor"
	"github.com/arefev/gophkeeper/internal/server/repository"
	"github.com/arefev/gophkeeper/internal/server/repository/testdb"
	"github.com/arefev/gophkeeper/internal/server/trm"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
)

type prepare struct {
	dbDSN     string
	dbConn    *postgresql.DB
	container *testdb.TestDBContainer
	srv       *grpc.Server
}

func NewPrepare() *prepare {
	return &prepare{}
}

func (p *prepare) runDB(ctx context.Context) error {
	db, err := testdb.New(ctx)
	if err != nil {
		return fmt.Errorf("testdb new failed: %w", err)
	}

	p.dbDSN = db.URI + "?sslmode=disable"
	p.container = db
	return nil
}

func (p *prepare) app() (*application.App, error) {
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("ENCRYPTION_SECRET", "thisis32bitlongpassphraseimusing")

	conf, err := config.NewConfig([]string{})
	if err != nil {
		return nil, fmt.Errorf("init config fail: %w", err)
	}

	l, err := logger.Build(conf.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("logger build fail: %w", err)
	}

	db, err := postgresql.NewDB(l).Connect(p.dbDSN)
	if err != nil {
		return nil, fmt.Errorf("db trm connect fail: %w", err)
	}

	err = p.migrationsUp(p.dbDSN)
	if err != nil {
		return nil, fmt.Errorf("run: migration up fail: %w", err)
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

	p.dbConn = db

	return app, nil
}

func (p *prepare) migrationsUp(dsn string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("migrations get dir failed: %w", err)
	}
	path := strings.Replace(dir, "internal/test", "cmd/server/db/migrations", -1)

	m, err := migrate.New("file://"+path, dsn)
	if err != nil {
		return fmt.Errorf("migrations instance fail: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrations up fail: %w", err)
	}

	return nil
}

func (p *prepare) server(app *application.App) error {
	listen, err := net.Listen("tcp", app.Conf.Address)
	if err != nil {
		return fmt.Errorf("server create listener failed: %w", err)
	}

	intr := interceptor.New(app)
	p.srv = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			intr.UnaryCheckToken(),
		),
		grpc.ChainStreamInterceptor(
			intr.StreamCheckToken(),
		),
	)
	proto.RegisterAuthServer(p.srv, handler.NewAuthHandler(app))
	proto.RegisterFileServer(p.srv, handler.NewFileHandler(app))
	proto.RegisterListServer(p.srv, handler.NewListHandler(app))
	go func() {
		if err := p.srv.Serve(listen); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return nil
}

func (p *prepare) Close(ctx context.Context) error {

	if p.srv != nil {
		p.srv.Stop()
	}

	if err := p.dbConn.Close(); err != nil {
		return err
	}

	if err := p.container.Close(ctx); err != nil {
		return err
	}

	return nil
}
