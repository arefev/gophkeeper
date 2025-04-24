package handler

import (
	"context"
	"fmt"

	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/application"
	"github.com/arefev/gophkeeper/internal/server/model"
)

type listHandler struct {
	proto.UnimplementedListServer
	app *application.App
}

func NewListHandler(app *application.App) *listHandler {
	return &listHandler{
		app: app,
	}
}

func (lh *listHandler) Get(
	ctx context.Context,
	in *proto.MetaListRequest,
) (*proto.MetaListResponse, error) {
	var list []model.Meta
	var err error
	err = lh.app.TrManager.Do(ctx, func(ctx context.Context) error {
		// TODO: подставлять ID авторизованного пользователя
		list, err = lh.app.Rep.Meta.Get(ctx, 1)
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

func (lh *listHandler) Delete(
	ctx context.Context,
	in *proto.MetaDeleteRequest,
) (*proto.MetaDeleteResponse, error) {
	var err error
	err = lh.app.TrManager.Do(ctx, func(ctx context.Context) error {
		// TODO: подставлять ID авторизованного пользователя
		err = lh.app.Rep.Meta.DeleteByUuid(ctx, in.GetUuid(), 1)
		if err != nil {
			return fmt.Errorf("run: meta delete failed: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("run: do transaction failed: %w", err)
	}

	return &proto.MetaDeleteResponse{}, nil
}
