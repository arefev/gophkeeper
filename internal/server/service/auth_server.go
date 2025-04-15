package service

import (
	"context"
	"fmt"

	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/application"
	"go.uber.org/zap"
)

type authServer struct {
	proto.UnimplementedAuthServer
	app *application.App
}

func NewAuthServer(app *application.App) *authServer {
	return &authServer{
		app: app,
	}
}

func (asrv *authServer) Register(
	ctx context.Context,
	in *proto.RegistrationRequest,
) (*proto.RegistrationResponse, error) {
	s := NewUserService(asrv.app)
	err := s.Create(ctx, in.User.GetLogin(), in.User.GetPassword())
	if err != nil {
		asrv.app.Log.Debug(
			"register user failed",
			zap.Error(err),
			zap.String("login", in.User.GetLogin()),
			zap.String("pwd", in.User.GetPassword()),
		)

		return &proto.RegistrationResponse{}, fmt.Errorf("register create user failed: %w", err)
	}

	asrv.app.Log.Debug(
		"register user success",
		zap.String("login", in.User.GetLogin()),
	)

	token, err := s.Authorize(ctx, in.User.GetLogin(), in.User.GetPassword())
	if err != nil {
		return &proto.RegistrationResponse{}, fmt.Errorf("register authorize user failed: %w", err)
	}

	return &proto.RegistrationResponse{Token: token.AccessToken}, nil
}

func (asrv *authServer) Login(
	ctx context.Context,
	in *proto.AuthorizationRequest,
) (*proto.AuthorizationResponse, error) {
	s := NewUserService(asrv.app)
	token, err := s.Authorize(ctx, in.User.GetLogin(), in.User.GetPassword())
	if err != nil {
		asrv.app.Log.Debug(
			"login user failed",
			zap.Error(err),
			zap.String("login", in.User.GetLogin()),
			zap.String("pwd", in.User.GetPassword()),
		)

		return &proto.AuthorizationResponse{}, fmt.Errorf("login authorize user failed: %w", err)
	}

	asrv.app.Log.Debug(
		"login user success",
		zap.String("login", in.User.GetLogin()),
	)

	return &proto.AuthorizationResponse{Token: token.AccessToken}, nil
}
