package server

import (
	"context"
	"fmt"

	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/application"
	"github.com/arefev/gophkeeper/internal/server/model"
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
	var list []model.Meta
	var err error
	err = ls.app.TrManager.Do(ctx, func(ctx context.Context) error {
		list, err = ls.app.Rep.Meta.Get(ctx, 2)
		if err != nil {
			return fmt.Errorf("run: meta get failed: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("run: do transaction failed: %w", err)
	}

	respList := []*proto.MetaList{}

	for _, item := range list {
		uuid := item.Uuid.String()
		t := item.Type.String()
		date := item.CreatedAt.Format("02.01.2006 15:04:05")
		meta := &proto.MetaList{
			Uuid:      &uuid,
			Type:      &t,
			Name:      &item.Name,
			FileName:  &item.File.Name,
			CreatedAt: &date,
		}
		respList = append(respList, meta)
	}

	resp := &proto.MetaListResponse{MetaList: respList}

	return resp, nil
}
