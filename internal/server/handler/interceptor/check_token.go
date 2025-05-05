package interceptor

import (
	"context"
	"errors"
	"fmt"

	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/handler/middleware"
	"github.com/arefev/gophkeeper/internal/server/model"
	"github.com/arefev/gophkeeper/internal/server/service"
	"github.com/arefev/gophkeeper/internal/server/service/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// StreamCheckToken is a grpc stream interceptor for checking authorization token
// return:
//
//	grpc.StreamServerInterceptor
func (i *interceptor) StreamCheckToken() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, err := i.checkToken(ss.Context())
		if err != nil {
			i.app.Log.Debug("stream check token failed", zap.Error(err))
			return status.Error(codes.Unauthenticated, "invalid token")
		}

		wrapped := middleware.WrapServerStream(ss)
		wrapped.WrappedContext = ctx

		return handler(srv, wrapped)
	}
}

// StreamCheckToken is a grpc unary interceptor for checking authorization token
// return:
//
//	grpc.StreamServerInterceptor
func (i *interceptor) UnaryCheckToken() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if _, ok := info.Server.(proto.AuthServer); ok {
			return handler(ctx, req)
		}

		ctx, err := i.checkToken(ctx)
		if err != nil {
			i.app.Log.Debug("unary check token failed", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		return handler(ctx, req)
	}
}

func (i *interceptor) checkToken(ctx context.Context) (context.Context, error) {
	var token string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing context")
	}

	values := md.Get("token")
	if len(values) > 0 {
		token = values[0]
	}

	if token == "" {
		return nil, errors.New("missing token")
	}

	login, err := jwt.NewToken(i.app.Conf.TokenSecret).Parse(token).GetLogin()
	if err != nil {
		return nil, fmt.Errorf("parse token failed: %w", err)
	}

	user, err := i.findUserByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("get user from context failed: %w", err)
	}

	ctx = context.WithValue(ctx, model.User{}, user)
	return ctx, nil
}

func (i *interceptor) findUserByLogin(ctx context.Context, login string) (*model.User, error) {
	var user *model.User

	user, err := service.NewUserService(i.app).GetUser(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("get user fail: %w", err)
	}

	return user, nil
}
