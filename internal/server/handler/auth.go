package handler

import (
	"context"
	"fmt"

	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/application"
	"github.com/arefev/gophkeeper/internal/server/service"
	"go.uber.org/zap"
)

type authHandler struct {
	proto.UnimplementedAuthServer
	app *application.App
}

func NewAuthHandler(app *application.App) *authHandler {
	return &authHandler{
		app: app,
	}
}

func (ah *authHandler) Register(
	ctx context.Context,
	in *proto.RegistrationRequest,
) (*proto.RegistrationResponse, error) {
	s := service.NewUserService(ah.app)
	err := s.Create(ctx, in.GetUser().GetLogin(), in.GetUser().GetPassword())
	if err != nil {
		return &proto.RegistrationResponse{}, fmt.Errorf("register create user failed: %w", err)
	}

	ah.app.Log.Debug(
		"register user success",
		zap.String("login", in.GetUser().GetLogin()),
	)

	token, err := s.Authorize(ctx, in.GetUser().GetLogin(), in.GetUser().GetPassword())
	if err != nil {
		return &proto.RegistrationResponse{}, fmt.Errorf("register authorize user failed: %w", err)
	}

	return &proto.RegistrationResponse{Token: &token.AccessToken}, nil
}

func (ah *authHandler) Login(
	ctx context.Context,
	in *proto.AuthorizationRequest,
) (*proto.AuthorizationResponse, error) {
	s := service.NewUserService(ah.app)
	token, err := s.Authorize(ctx, in.GetUser().GetLogin(), in.GetUser().GetPassword())
	if err != nil {
		return &proto.AuthorizationResponse{}, fmt.Errorf("login authorize user failed: %w", err)
	}

	ah.app.Log.Debug(
		"login user success",
		zap.String("login", in.GetUser().GetLogin()),
	)

	return &proto.AuthorizationResponse{Token: &token.AccessToken}, nil
}
