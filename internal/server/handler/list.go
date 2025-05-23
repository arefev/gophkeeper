package handler

import (
	"context"
	"fmt"

	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/application"
	"github.com/arefev/gophkeeper/internal/server/model"
	"github.com/arefev/gophkeeper/internal/server/service"
)

type listHandler struct {
	proto.UnimplementedListServer
	app *application.App
}

// NewListHandler create and return pointer on new scruct listHandler
//
//	app - pointer on struct application.App
func NewListHandler(app *application.App) *listHandler {
	return &listHandler{
		app: app,
	}
}

// Get is a grpc handler for get of list meta data
// params:
//
//	ctx context.Context
//	in *proto.MetaListRequest - request with data for get list
//
// return:
//
//	*proto.MetaListResponse - has data with list of meta
//	error
func (lh *listHandler) Get(
	ctx context.Context,
	in *proto.MetaListRequest,
) (*proto.MetaListResponse, error) {
	var list []model.Meta
	var err error

	user, err := service.NewUserService(lh.app).Authorized(ctx)
	if err != nil {
		return nil, service.ErrAuthUserNotFound
	}

	err = lh.app.TrManager.Do(ctx, func(ctx context.Context) error {
		list, err = lh.app.Rep.Meta.Get(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("meta get failed: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("do transaction failed: %w", err)
	}

	respList := []*proto.MetaList{}

	for i := range list {
		uuid := list[i].UUID.String()
		t := list[i].Type.String()
		date := list[i].CreatedAt.Format("02.01.2006 15:04:05")

		meta := proto.MetaList_builder{
			Uuid:      &uuid,
			Type:      &t,
			Name:      &list[i].Name,
			FileName:  &list[i].File.Name,
			CreatedAt: &date,
		}.Build()

		respList = append(respList, meta)
	}

	resp := proto.MetaListResponse_builder{MetaList: respList}.Build()
	return resp, nil
}

// Delete is a grpc handler for delete item from list meta data
// params:
//
//	ctx context.Context
//	in *proto.MetaDeleteRequest - request with data for delete item, etc. uuid
//
// return:
//
//	*proto.MetaDeleteResponse - grpc response
//	error
func (lh *listHandler) Delete(
	ctx context.Context,
	in *proto.MetaDeleteRequest,
) (*proto.MetaDeleteResponse, error) {
	var err error
	user, err := service.NewUserService(lh.app).Authorized(ctx)
	if err != nil {
		return nil, service.ErrAuthUserNotFound
	}

	err = lh.app.TrManager.Do(ctx, func(ctx context.Context) error {
		err = lh.app.Rep.Meta.DeleteByUUID(ctx, in.GetUuid(), user.ID)
		if err != nil {
			return fmt.Errorf("meta delete failed: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("do transaction failed: %w", err)
	}

	return &proto.MetaDeleteResponse{}, nil
}
