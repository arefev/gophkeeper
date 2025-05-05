package application

import (
	"context"

	"github.com/arefev/gophkeeper/internal/server/config"
	"github.com/arefev/gophkeeper/internal/server/model"
	"github.com/arefev/gophkeeper/internal/server/trm"
	"go.uber.org/zap"
)

type UserRepo interface {
	Exists(ctx context.Context, login string) bool
	FindByLogin(ctx context.Context, login string) (*model.User, bool)
	Create(ctx context.Context, login, password string) error
}

type MetaRepo interface {
	Create(ctx context.Context, meta *model.Meta) error
	FindByUUID(ctx context.Context, uuid string, userID int) (*model.Meta, error)
	Get(ctx context.Context, userID int) ([]model.Meta, error)
	DeleteByUUID(ctx context.Context, uuid string, userID int) error
}

type TrManager interface {
	Do(ctx context.Context, action trm.TrAction) error
}

type App struct {
	Rep       Repository
	TrManager TrManager
	Log       *zap.Logger
	Conf      *config.Config
}

type Repository struct {
	User UserRepo
	Meta MetaRepo
}
