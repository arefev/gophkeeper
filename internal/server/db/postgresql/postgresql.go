package postgresql

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type DB struct {
	conn *sqlx.DB
	log  *zap.Logger
}

func NewDB(log *zap.Logger) *DB {
	return &DB{
		log: log,
	}
}

func (db *DB) Connect(dsn string) (*DB, error) {
	conn, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("db connect fail: %w", err)
	}

	db.conn = conn

	return db, nil
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
