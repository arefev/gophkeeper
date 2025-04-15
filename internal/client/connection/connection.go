package connection

import (
	"context"
	"fmt"

	"github.com/arefev/gophkeeper/internal/proto"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CheckAuthFail bool

type grpcClient struct {
	conn  *grpc.ClientConn
	log   *zap.Logger
	token string
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

func (g *grpcClient) SetToken(t string) {
	if len(t) == 0 {
		return
	}

	g.token = t
}

func (g *grpcClient) CheckTokenCmd() tea.Msg {
	// TODO: возможно нужна проверка на актуальность токена
	if g.token == "" {
		return CheckAuthFail(true)
	}
	return nil
}

func (g *grpcClient) Register(ctx context.Context, login, pwd string) (string, error) {
	client := proto.NewRegistrationClient(g.conn)

	resp, err := client.Register(ctx, &proto.RegistrationRequest{
		User: &proto.User{
			Login:    login,
			Password: pwd,
		},
	})

	if err != nil {
		return "", fmt.Errorf("grpc Register failed: %w", err)
	}

	return resp.GetToken(), nil
}
