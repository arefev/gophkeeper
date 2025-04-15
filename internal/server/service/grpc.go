package service

import (
	"context"
	"fmt"

	"github.com/arefev/gophkeeper/internal/proto"
)

type GRPCServer struct {
	proto.UnimplementedRegistrationServer
}

func (gs *GRPCServer) Register(
	ctx context.Context,
	in *proto.RegistrationRequest,
) (*proto.RegistrationResponse, error) {
	fmt.Println("message recieved")
	return &proto.RegistrationResponse{}, nil
}
