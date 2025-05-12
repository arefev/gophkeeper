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

// NewAuthHandler create and return pointer on new scruct authHandler
//
//	app - pointer on struct application.App
func NewAuthHandler(app *application.App) *authHandler {
	return &authHandler{
		app: app,
	}
}

// Register is a grpc handler for register new user
// params:
//
//	ctx - context
//	in - pointer on struct proto.RegistrationRequest, has request data for user registration, etc. login & password
//
// return:
//
//	*proto.RegistrationResponse - grpc response
//	error
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

	resp := proto.RegistrationResponse_builder{Token: &token.AccessToken}.Build()
	return resp, nil
}

// Login is a grpc handler for authorization user
// params:
//
//	ctx - context
//	in - pointer on struct proto.AuthorizationRequest, has request data for user authorization, etc. login & password
//
// return:
//
//	*proto.AuthorizationResponse - grpc response
//	error
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

	resp := proto.AuthorizationResponse_builder{Token: &token.AccessToken}.Build()
	return resp, nil
}
