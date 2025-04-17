package connection

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
	client := proto.NewAuthClient(g.conn)

	resp, err := client.Register(ctx, &proto.RegistrationRequest{
		User: &proto.User{
			Login:    &login,
			Password: &pwd,
		},
	})

	if err != nil {
		return "", fmt.Errorf("grpc Register failed: %w", err)
	}

	return resp.GetToken(), nil
}

func (g *grpcClient) Login(ctx context.Context, login, pwd string) (string, error) {
	client := proto.NewAuthClient(g.conn)

	resp, err := client.Login(ctx, &proto.AuthorizationRequest{
		User: &proto.User{
			Login:    &login,
			Password: &pwd,
		},
	})

	if err != nil {
		return "", fmt.Errorf("grpc Login failed: %w", err)
	}

	return resp.GetToken(), nil
}

func (g *grpcClient) TextUpload(ctx context.Context, txt []byte, metaName, metaType string) error {
	client := proto.NewFileClient(g.conn)
	stream, err := client.Upload(ctx)
	if err != nil {
		return fmt.Errorf("grpc text upload stream failed: %w", err)
	}

	fileName := "test.txt"
	err = stream.Send(&proto.FileUploadRequest{
		Chunk: txt,
		Name:  &fileName,
		Meta: &proto.Meta{
			Name: &metaName,
			Type: &metaType,
		},
	})
	if err != nil {
		return fmt.Errorf("grpc text upload send failed: %w", err)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("grpc text upload close failed: %w", err)
	}
	g.log.Debug("text sent", zap.Uint32("bytes", res.GetSize()))

	return nil
}

func (g *grpcClient) FileUpload(ctx context.Context, path, metaName, metaType string) error {
	client := proto.NewFileClient(g.conn)
	stream, err := client.Upload(ctx)
	if err != nil {
		return fmt.Errorf("grpc creds upload stream failed: %w", err)
	}

	fileName := filepath.Base(path)

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	buf := make([]byte, 1024*1024) // 1МБ
	batchNumber := 1
	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		chunk := buf[:num]

		r := &proto.FileUploadRequest{
			Chunk: chunk,
			Name:  &fileName,
			Meta: &proto.Meta{
				Name: &metaName,
				Type: &metaType,
			},
		}
		if err := stream.Send(r); err != nil {
			return err
		}
		g.log.Debug("Sent - batch", zap.Int("number", batchNumber), zap.Int("size", len(chunk)))
		batchNumber += 1

	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	g.log.Debug("file uploaded", zap.Uint32("bytes", res.GetSize()))
	return nil
}
