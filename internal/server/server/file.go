package server

import (
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

	return nil
}
