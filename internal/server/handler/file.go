package handler

import (
	"context"
	"fmt"

	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/application"
	"github.com/arefev/gophkeeper/internal/server/model"
	"github.com/arefev/gophkeeper/internal/server/service"
	"go.uber.org/zap"
)

type fileHandler struct {
	proto.UnimplementedFileServer
	app *application.App
}

func NewFileHandler(app *application.App) *fileHandler {
	return &fileHandler{
		app: app,
	}
}

func (fh *fileHandler) Upload(stream proto.File_UploadServer) error {
	user, err := service.NewUserService(fh.app).Authorized(stream.Context())
	if err != nil {
		return service.ErrAuthUserNotFound
	}

	storage := service.NewStorageService(fh.app)
	err = storage.Upload(user.ID, stream)
	if err != nil {
		fh.app.Log.Debug(
			"file upload failed",
			zap.Error(err),
		)
		return fmt.Errorf("file upload failed: %w", err)
	}

	err = storage.Save()
	if err != nil {
		fh.app.Log.Debug(
			"file data save failed",
			zap.Error(err),
		)
		return fmt.Errorf("file data save failed: %w", err)
	}

	err = storage.Remove()
	if err != nil {
		fh.app.Log.Debug(
			"file remove failed",
			zap.Error(err),
		)
		return fmt.Errorf("file remove failed: %w", err)
	}

	return nil
}

func (fh *fileHandler) Download(
	req *proto.FileDownloadRequest,
	stream proto.File_DownloadServer,
) error {
	var data []byte
	var meta *model.Meta
	var err error

	user, err := service.NewUserService(fh.app).Authorized(stream.Context())
	if err != nil {
		return service.ErrAuthUserNotFound
	}

	err = fh.app.TrManager.Do(stream.Context(), func(ctx context.Context) error {
		meta, err = fh.app.Rep.Meta.FindByUUID(ctx, req.GetUuid(), user.ID)
		if err != nil {
			return fmt.Errorf("meta find by uuid failed: %w", err)
		}

		es := service.NewEncryptionService(fh.app)
		data, err = es.Decrypt(meta.File.Data)
		if err != nil {
			return fmt.Errorf("decrypt data failed: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("do transaction failed: %w", err)
	}

	read := 0
	max := fh.app.Conf.ChunkSize
	size := len(data)

	for read < size {
		next := min(read+max, size)
		chunk := data[read:next]

		err := stream.Send(&proto.FileDownloadResponse{Chunk: chunk, Name: &meta.File.Name})
		if err != nil {
			return fmt.Errorf("download stream failed: %w", err)
		}

		read = next
	}

	return nil
}
