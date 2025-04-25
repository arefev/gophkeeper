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

func (m *Meta) FindByUUID(ctx context.Context, uuid string, userID int) (*model.Meta, error) {
	ctx, cancel := context.WithTimeout(ctx, timeCancel)
	defer cancel()

	q := `
		SELECT 
			m.id, m.uuid, m.type, m.name, m.user_id, m.created_at, m.updated_at,
			f.id as f_id, f.meta_id as f_meta_id, f.name as f_name,
			f.data as f_data, f.created_at as f_created_at 
		FROM meta as m
		JOIN files as f ON m.id = f.meta_id
		WHERE m.uuid = :uuid AND m.user_id = :user_id
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
	arg := map[string]any{"uuid": uuid, "user_id": userID}
	row := stmt.QueryRow(arg)
	err = row.Scan(
		&meta.ID,
		&meta.UUID,
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

func (m *Meta) DeleteByUUID(ctx context.Context, uuid string, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, timeCancel)
	defer cancel()

	q := "DELETE FROM meta WHERE user_id = :user_id AND uuid = :uuid"
	args := map[string]interface{}{
		"uuid":    uuid,
		"user_id": userID,
	}

	if err := m.execWithArgs(ctx, args, q); err != nil {
		return fmt.Errorf("meta delete fail: %w", err)
	}

	return nil
}

func (m *Meta) Get(ctx context.Context, userID int) ([]model.Meta, error) {
	ctx, cancel := context.WithTimeout(ctx, timeCancel)
	defer cancel()

	q := `
		SELECT 
			m.id, m.uuid, m.type, m.name, m.user_id, m.created_at, m.updated_at,
			f.id as f_id, f.meta_id as f_meta_id, f.name as f_name,
			f.created_at as f_created_at 
		FROM meta as m
		JOIN files as f ON m.id = f.meta_id
		WHERE m.user_id = :user_id
		ORDER BY m.created_at DESC
	`

	stmt, err := m.prepare(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("meta get: prepare failed: %w", err)
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			m.log.Warn("meta get: stmt close failed", zap.Error(err))
		}
	}()

	list := []model.Meta{}
	arg := map[string]any{"user_id": userID}
	rows, err := stmt.Queryx(arg)
	if err != nil {
		return nil, fmt.Errorf("meta get: query failed: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			m.log.Warn("meta get: rows close failed", zap.Error(err))
		}
	}()

	for rows.Next() {
		meta := model.Meta{}
		err = rows.Scan(
			&meta.ID,
			&meta.UUID,
			&meta.Type,
			&meta.Name,
			&meta.UserID,
			&meta.CreatedAt,
			&meta.UpdatedAt,
			&meta.File.ID,
			&meta.File.MetaID,
			&meta.File.Name,
			&meta.File.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("meta get: scan failed: %w", err)
		}
		list = append(list, meta)
	}

	return list, nil
}
