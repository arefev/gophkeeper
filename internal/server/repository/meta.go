package repository

import (
	"context"
	"fmt"

	"github.com/arefev/gophkeeper/internal/server/model"
	"go.uber.org/zap"
)

type Meta struct {
	log *zap.Logger
	*Base
}

func NewMeta(tr TxGetter, log *zap.Logger) *Meta {
	return &Meta{
		log:  log,
		Base: NewBase(tr, log),
	}
}

func (m *Meta) Create(ctx context.Context, meta *model.Meta) error {
	ctx, cancel := context.WithTimeout(ctx, timeCancel)
	defer cancel()

	query := `
		WITH inserted AS (
			INSERT INTO meta(user_id, type, name) VALUES(:user_id, :type, :name) RETURNING id
		)
		INSERT INTO files(name, data, meta_id) VALUES(:file_name, :file_data, (SELECT id FROM inserted))
	`

	args := map[string]any{
		"user_id":   meta.UserID,
		"type":      meta.Type,
		"name":      meta.Name,
		"file_name": meta.File.Name,
		"file_data": meta.File.Data,
	}

	err := m.execWithArgs(ctx, args, query)
	if err != nil {
		return fmt.Errorf("meta create failed: %w", err)
	}

	return nil
}

func (m *Meta) Find(ctx context.Context, id int) (*model.Meta, error) {
	ctx, cancel := context.WithTimeout(ctx, timeCancel)
	defer cancel()

	q := `
		SELECT 
			m.id, m.uuid, m.type, m.name, m.user_id, m.created_at, m.updated_at,
			f.id as f_id, f.meta_id as f_meta_id, f.name as f_name,
			f.data as f_data, f.created_at as f_created_at 
		FROM meta as m
		JOIN files as f ON m.id = f.meta_id
		WHERE m.id = :id
	`

	stmt, err := m.prepare(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("exec with args: prepare fail: %w", err)
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			m.log.Warn("create with args: stmt close fail", zap.Error(err))
		}
	}()

	meta := &model.Meta{}
	arg := map[string]interface{}{"id": id}
	row := stmt.QueryRow(arg)
	err = row.Scan(
		&meta.ID,
		&meta.Uuid,
		&meta.Type,
		&meta.Name,
		&meta.UserID,
		&meta.CreatedAt,
		&meta.UpdatedAt,
		&meta.File.ID,
		&meta.File.MetaID,
		&meta.File.Name,
		&meta.File.Data,
		&meta.File.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return meta, nil
}
