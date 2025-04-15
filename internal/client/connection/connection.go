package connection

import (
	"context"
	"fmt"

	"github.com/arefev/gophkeeper/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcClient struct {
	conn *grpc.ClientConn
	log  *zap.Logger
}

func NewGRPCClient(l *zap.Logger) *grpcClient {
	return &grpcClient{
		log: l,
	}
}

func (g *grpcClient) Connect(address string) error {
	var err error
	g.conn, err = grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("grpc connect failed: %w", err)
	}
	return nil
}

func (g *grpcClient) Close() error {
	if err := g.conn.Close(); err != nil {
		return fmt.Errorf("grpc connection close failed: %w", err)
	}
	return nil
}

func (g *grpcClient) Register(ctx context.Context, login, pwd string) error {
	client := proto.NewRegistrationClient(g.conn)

	_, err := client.Register(ctx, &proto.RegistrationRequest{
		User: &proto.User{
			Login: login,
			Password: pwd,
		},
	})

	if err != nil {
		return fmt.Errorf("grpc Register failed: %w", err)
	}

	return nil
}
