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
}
