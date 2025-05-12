// Package using for:
//
//	 create connection to grpc server
//		send request to register user
//		send request to authorize user
//		send request to upload on server text or file
//		send request to get list of meta data
//		send request to delete item from meta data
//		send reqeust to download file from server
package connection

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/arefev/gophkeeper/internal/client/service"
	"github.com/arefev/gophkeeper/internal/client/tui/model"
	"github.com/arefev/gophkeeper/internal/proto"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type CheckAuthFail bool

type grpcClient struct {
	conn      *grpc.ClientConn
	log       *zap.Logger
	token     string
	chunkSize int
}

// NewGRPCClient create object grpcClient and return pointer.
func NewGRPCClient(chunkSize int, l *zap.Logger) *grpcClient {
	return &grpcClient{
		chunkSize: chunkSize,
		log:       l,
	}
}

// Connnect create connection to grpc server.
func (g *grpcClient) Connect(address string) error {
	var err error
	g.conn, err = grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("grpc connect failed: %w", err)
	}
	return nil
}

// Close closes connection.
func (g *grpcClient) Close() error {
	if err := g.conn.Close(); err != nil {
		return fmt.Errorf("grpc connection close failed: %w", err)
	}
	return nil
}

// SetToken sets authorization token.
func (g *grpcClient) SetToken(t string) {
	if len(t) == 0 {
		return
	}

	g.token = t
}

// CheckTokenCmd check token on exist.
func (g *grpcClient) CheckTokenCmd() tea.Msg {
	if g.token == "" {
		return CheckAuthFail(true)
	}
	return nil
}

// Register send request to create new user by login and pwd.
// Return authorization token and error.
func (g *grpcClient) Register(ctx context.Context, login, pwd string) (string, error) {
	client := proto.NewAuthClient(g.conn)

	resp, err := client.Register(ctx, proto.RegistrationRequest_builder{
		User: proto.User_builder{
			Login:    &login,
			Password: &pwd,
		}.Build(),
	}.Build())

	if err != nil {
		return "", fmt.Errorf("grpc Register failed: %w", err)
	}

	return resp.GetToken(), nil
}

// Register send request to authorize user by login and pwd.
// Return authorization token and error.
func (g *grpcClient) Login(ctx context.Context, login, pwd string) (string, error) {
	client := proto.NewAuthClient(g.conn)

	resp, err := client.Login(ctx, proto.AuthorizationRequest_builder{
		User: proto.User_builder{
			Login:    &login,
			Password: &pwd,
		}.Build(),
	}.Build())

	if err != nil {
		return "", fmt.Errorf("grpc Login failed: %w", err)
	}

	return resp.GetToken(), nil
}

// TextUpload send request to upload on server text information.
// Params:
//
//	txt []byte - information in bytes
//	metaName string - name of meta data
//	metaType string - type ofr meta data [creds, bank, etc...]
//
// Return error.
func (g *grpcClient) TextUpload(ctx context.Context, txt []byte, metaName, metaType string) error {
	md := metadata.New(map[string]string{"token": g.token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	client := proto.NewFileClient(g.conn)
	stream, err := client.Upload(ctx)
	if err != nil {
		return fmt.Errorf("grpc text upload stream failed: %w", err)
	}

	fileName := metaType + ".txt"
	err = stream.Send(proto.FileUploadRequest_builder{
		Chunk: txt,
		Name:  &fileName,
		Meta: proto.Meta_builder{
			Name: &metaName,
			Type: &metaType,
		}.Build(),
	}.Build())
	if err != nil {
		return fmt.Errorf("grpc text upload send failed: %w", err)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		g.log.Sugar().Infof("stream err %+v", status.Code(err))
		return fmt.Errorf("grpc text upload close failed: %w", err)
	}
	g.log.Debug("text sent", zap.Uint32("bytes", res.GetSize()))

	return nil
}

// FileUpload send request to upload on server bynary file.
// Params:
//
//	path string - path to file for upload
//	metaName string - name of meta data
//	metaType string - type of meta data
//
// Return error.
func (g *grpcClient) FileUpload(ctx context.Context, path, metaName, metaType string) error {
	md := metadata.New(map[string]string{"token": g.token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	client := proto.NewFileClient(g.conn)
	stream, err := client.Upload(ctx)
	if err != nil {
		return fmt.Errorf("grpc file upload stream failed: %w", err)
	}

	fileName := filepath.Base(path)

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("grpc file upload open failed: %w", err)
	}
	buf := make([]byte, g.chunkSize)
	batchNumber := 1
	for {
		num, err := file.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("grpc file upload read failed: %w", err)
		}
		chunk := buf[:num]

		r := proto.FileUploadRequest_builder{
			Chunk: chunk,
			Name:  &fileName,
			Meta: proto.Meta_builder{
				Name: &metaName,
				Type: &metaType,
			}.Build(),
		}.Build()
		if err := stream.Send(r); err != nil {
			return fmt.Errorf("grpc file upload stream send failed: %w", err)
		}
		g.log.Debug("Sent - batch", zap.Int("number", batchNumber), zap.Int("size", len(chunk)))
		batchNumber++
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("grpc file upload stream close failed: %w", err)
	}

	g.log.Debug("file uploaded", zap.Uint32("bytes", res.GetSize()))
	return nil
}

// GetList send request to get list meta data for user.
// Return pointer on slice with meta data, and error.
func (g *grpcClient) GetList(ctx context.Context) (*[]model.MetaListData, error) {
	md := metadata.New(map[string]string{"token": g.token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	client := proto.NewListClient(g.conn)
	resp, err := client.Get(ctx, &proto.MetaListRequest{})

	if err != nil {
		return nil, fmt.Errorf("grpc get list failed: %w", err)
	}

	list := []model.MetaListData{}
	for _, item := range resp.GetMetaList() {
		list = append(list, model.MetaListData{
			UUID: item.GetUuid(),
			Type: item.GetType(),
			Name: item.GetName(),
			File: item.GetFileName(),
			Date: item.GetCreatedAt(),
		})
	}

	return &list, nil
}

// Delete send request to delete by uuid item from meta data for user.
// Return error.
func (g *grpcClient) Delete(ctx context.Context, uuid string) error {
	md := metadata.New(map[string]string{"token": g.token})
	ctx = metadata.NewOutgoingContext(ctx, md)
	client := proto.NewListClient(g.conn)

	_, err := client.Delete(ctx, proto.MetaDeleteRequest_builder{Uuid: &uuid}.Build())
	if err != nil {
		return fmt.Errorf("grpc FileDelete failed: %w", err)
	}

	return nil
}

// FileDownload send request to download file (text or bynary) from server.
// Return path of downloaded file or error.
func (g *grpcClient) FileDownload(ctx context.Context, uuid string) (string, error) {
	md := metadata.New(map[string]string{"token": g.token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	client := proto.NewFileClient(g.conn)

	req := proto.FileDownloadRequest_builder{Uuid: &uuid}.Build()
	stream, err := client.Download(ctx, req)

	if err != nil {
		return "", fmt.Errorf("grpc file download stream failed: %w", err)
	}

	file := service.NewFile()
	var fileSize uint32 = 0

	defer func() {
		if err := file.Output.Close(); err != nil {
			g.log.Error("file download close failed", zap.Error(err))
		}
	}()
	for {
		req, err := stream.Recv()
		if file.Path == "" {
			dir, err := os.Getwd()
			if err != nil {
				return "", fmt.Errorf("file download get dir failed: %w", err)
			}

			if err = file.SetFile(req.GetName(), dir); err != nil {
				return "", fmt.Errorf("file download create file failed: %w", err)
			}
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			g.log.Debug("file download stream recv failed", zap.Error(err))
			return "", fmt.Errorf("file download stream recv failed: %w", err)
		}
		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))
		g.log.Debug("received a chunk with size", zap.Uint32("size", fileSize))
		if err := file.Write(chunk); err != nil {
			g.log.Debug("file download write failed", zap.Error(err))
			return "", fmt.Errorf("file download write failed: %w", err)
		}
	}

	g.log.Debug("file downloaded", zap.Uint32("size", fileSize))

	return file.Path, nil
}
