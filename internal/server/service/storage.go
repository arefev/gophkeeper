package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/application"
	"github.com/arefev/gophkeeper/internal/server/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type storageService struct {
	app  *application.App
	file *FileService
	meta *model.Meta
}

func NewStorageService(app *application.App) *storageService {
	return &storageService{
		app:  app,
		file: NewFile(),
	}
}

func (s *storageService) Upload(userID int, stream proto.File_UploadServer) error {
	var fileSize uint32 = 0

	defer func() {
		if err := s.file.Output.Close(); err != nil {
			s.app.Log.Error("file upload close failed", zap.Error(err))
		}
	}()
	for {
		req, err := stream.Recv()
		if s.file.Path == "" {
			err := s.file.SetFile(req.GetName(), "./storage/"+uuid.NewString())
			if err != nil {
				return fmt.Errorf("file upload set file failed %w", err)
			}
			s.setMeta(userID, req.GetMeta().GetName(), req.GetMeta().GetType(), req.GetName())
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("file upload stream recv failed: %w", err)
		}
		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))
		s.app.Log.Debug("received a chunk with size", zap.Uint32("size", fileSize))
		if err := s.file.Write(chunk); err != nil {
			return fmt.Errorf("file upload write failed: %w", err)
		}
	}

	s.app.Log.Debug("file uploaded", zap.Uint32("size", fileSize))

	resp := proto.FileUploadResponse_builder{}.Build()
	resp.SetSize(fileSize)
	if err := stream.SendAndClose(resp); err != nil {
		return fmt.Errorf("file upload SendAndClose failed: %w", err)
	}

	return nil
}

func (s *storageService) setMeta(userID int, mName, mtype, fName string) {
	s.meta = &model.Meta{
		UserID: userID,
		Name:   mName,
		Type:   model.MetaTypeFromString(mtype),
		File: model.File{
			Name: fName,
		},
	}
}

func (s *storageService) Save() error {
	ctx := context.Background()
	var err error
	es := NewEncryptionService(s.app)
	s.meta.File.Data, err = s.file.ReadAll()
	if err != nil {
		return fmt.Errorf("storage service: readall from file failed: %w", err)
	}

	s.meta.File.Data, err = es.Encrypt(s.meta.File.Data)
	if err != nil {
		return fmt.Errorf("storage service: readall from file failed: %w", err)
	}

	err = s.app.TrManager.Do(ctx, func(ctx context.Context) error {
		err = s.app.Rep.Meta.Create(ctx, s.meta)
		if err != nil {
			return fmt.Errorf("storage service: meta create failed: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("storage service: do transaction failed: %w", err)
	}

	return nil
}

func (s *storageService) Remove() error {
	dir := filepath.Dir(s.file.Path)

	err := os.Remove(s.file.Path)
	if err != nil {
		return fmt.Errorf("storage service: remove file failed: %w", err)
	}

	err = os.Remove(dir)
	if err != nil {
		return fmt.Errorf("storage service: remove dir failed: %w", err)
	}

	return nil
}
