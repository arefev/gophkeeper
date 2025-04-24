package server

import (
	"context"
	"fmt"

	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/application"
	"github.com/arefev/gophkeeper/internal/server/service"
	"go.uber.org/zap"
)

type fileServer struct {
	proto.UnimplementedFileServer
	app *application.App
}

func NewFileServer(app *application.App) *fileServer {
	return &fileServer{
		app: app,
	}
}

func (fs *fileServer) Upload(stream proto.File_UploadServer) error {
	storage := service.NewStorageService(fs.app)
	err := storage.Upload(stream)
	if err != nil {
		fs.app.Log.Debug(
			"file upload failed",
			zap.Error(err),
		)
		return fmt.Errorf("file upload failed: %w", err)
	}

	err = storage.Save()
	if err != nil {
		fs.app.Log.Debug(
			"file data save failed",
			zap.Error(err),
		)
		return fmt.Errorf("file data save failed: %w", err)
	}

	err = storage.Remove()
	if err != nil {
		fs.app.Log.Debug(
			"file remove failed",
			zap.Error(err),
		)
		return fmt.Errorf("file remove failed: %w", err)
	}

	return nil
}

func (fs *fileServer) Download(
	req *proto.FileDownloadRequest,
	stream proto.File_DownloadServer,
) error {
	// TODO: получать userID из контекста
	const userID = 2

	err := fs.app.TrManager.Do(stream.Context(), func(ctx context.Context) error {
		meta, err := fs.app.Rep.Meta.FindByUuid(ctx, req.GetUuid(), userID)
		if err != nil {
			return fmt.Errorf("run: meta get failed: %w", err)
		}
		// l.Sugar().Infof("meta %+v", meta)

		es := service.NewEncryptionService(fs.app)
		_, err = es.Decrypt(meta.File.Data)
		if err != nil {
			return fmt.Errorf("run: decrypt data failed: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("run: do transaction failed: %w", err)
	}

	return nil
}
