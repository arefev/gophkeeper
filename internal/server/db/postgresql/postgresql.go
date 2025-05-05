package postgresql

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DB struct {
	conn *sqlx.DB
	log  *zap.Logger
	dsn  string
}

func NewDB(log *zap.Logger) *DB {
	return &DB{
		log: log,
	}
}

func (db *DB) Connect() error {
	conn, err := sqlx.Connect("pgx", db.dsn)
	if err != nil {
		return fmt.Errorf("db connect fail: %w", err)
	}

	db.conn = conn

	return nil
}

func (db *DB) Connection() *sqlx.DB {
	return db.conn
}

func (db *DB) Close() error {
	if err := db.conn.Close(); err != nil {
		return fmt.Errorf("db close fail: %w", err)
	}

	db.log.Info("db connection closed")
	return nil
}

func (db *DB) MigrationsUp() error {
	const migDir = "internal/server/db/postgresql/migrations"

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("migrations get dir failed: %w", err)
	}

	re := regexp.MustCompile(`internal.*`)
	f := re.FindString(dir)

	var path string
	if f != "" {
		path = strings.ReplaceAll(dir, f, migDir)
	} else {
		path = dir + "/" + migDir
	}

	m, err := migrate.New("file://"+path, db.dsn)
	if err != nil {
		return fmt.Errorf("migrations instance fail: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrations up fail: %w", err)
	}

	return nil
}

func (db *DB) DSNFromCreds(host, port, name, login, password string) *DB {
	host = net.JoinHostPort(host, port)
	db.dsn = fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		login,
		password,
		host,
		name,
	)

	return db
}

func (db *DB) DSN(dsn string) *DB {
	db.dsn = dsn
	return db
}
