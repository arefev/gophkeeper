package interceptor

import (
	"context"
	"fmt"

	"github.com/arefev/gophkeeper/internal/server/model"
	"github.com/arefev/gophkeeper/internal/server/service"
	"github.com/arefev/gophkeeper/internal/server/service/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (m *middleware) StreamCheckToken(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	var token string
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Error(codes.Unauthenticated, "missing context")
	}

	values := md.Get("token")
	if len(values) > 0 {
		token = values[0]
	}

	if token == "" {
		return status.Error(codes.Unauthenticated, "missing token")
	}

	login, err := jwt.NewToken(m.app.Conf.TokenSecret).Parse(token).GetLogin()
	if err != nil {
		m.app.Log.Debug("get login fail", zap.Error(err))
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	user, err := m.getUser(ss.Context(), login)
	if err != nil {
		m.app.Log.Debug("get user fail", zap.Error(err))
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	m.app.Log.Debug("incoming context", zap.String("token", token), zap.String("login", login), zap.Any("user", user))
	// ctx := context.WithValue(ss.Context(), model.User{}, user)
	// ctx = metadata.AppendToOutgoingContext(ctx)

	err = handler(srv, ss)
	return err
}

func (m *middleware) getUser(ctx context.Context, login string) (*model.User, error) {
	var user *model.User

	user, err := service.NewUserService(m.app).GetUser(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("get user fail: %w", err)
	}

	return user, nil
}
