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

func (m *Meta) Create(ctx context.Context, meta *model.Meta, file *model.File) error {
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
		"file_name": file.Name,
		"file_data": file.Data,
	}

	err := m.execWithArgs(ctx, args, query)
	if err != nil {
		return fmt.Errorf("meta create failed: %w", err)
	}

	return nil
}
