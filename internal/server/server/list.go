package server

import (
	"context"

	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/application"
)

type listServer struct {
	proto.UnimplementedListServer
	app *application.App
}

func NewListServer(app *application.App) *listServer {
	return &listServer{
		app: app,
	}
}

func (ls *listServer) Get(
	ctx context.Context,
	in *proto.MetaListRequest,
) (*proto.MetaListResponse, error) {
	return &proto.MetaListResponse{}, nil
}
