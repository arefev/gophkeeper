package service

import (
	"context"
	"fmt"

	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/application"
	"go.uber.org/zap"
)

type regServer struct {
	proto.UnimplementedRegistrationServer
	app *application.App
}

func NewRegServer(app *application.App) *regServer {
	return &regServer{
		app: app,
	}
}

func (rs *regServer) Register(
	ctx context.Context,
	in *proto.RegistrationRequest,
) (*proto.RegistrationResponse, error) {
	s := NewUserService(rs.app)
	err := s.Create(ctx, in.User.GetLogin(), in.User.GetPassword())
	if err != nil {
		rs.app.Log.Debug(
			"register user failed",
			zap.Error(err),
			zap.String("login", in.User.GetLogin()),
			zap.String("pwd", in.User.GetPassword()),
		)

		return &proto.RegistrationResponse{}, fmt.Errorf("register create user failed: %w", err)
	}

	rs.app.Log.Debug(
		"register user success",
		zap.String("login", in.User.GetLogin()),
	)

	token, err := s.Authorize(ctx, in.User.GetLogin(), in.User.GetPassword())
	if err != nil {
		return &proto.RegistrationResponse{}, fmt.Errorf("register authorize user failed: %w", err)
	}

	return &proto.RegistrationResponse{Token: token.AccessToken}, nil
}
