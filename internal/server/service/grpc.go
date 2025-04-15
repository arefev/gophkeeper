package service

import (
	"context"

	"github.com/arefev/gophkeeper/internal/proto"
)

type GRPCServer struct {
	proto.UnimplementedRegistrationServer
}

func (gs *GRPCServer) Register(
	ctx context.Context,
	in *proto.RegistrationRequest,
) (*proto.RegistrationResponse, error) {

	return &proto.RegistrationResponse{}, nil
}
