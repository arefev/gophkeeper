package server

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/arefev/gophkeeper/internal/proto"
	"github.com/arefev/gophkeeper/internal/server/application"
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

func (c *fileServer) Upload(stream proto.File_UploadServer) error {
	file := NewFile()
	var fileSize uint32
	fileSize = 0
	defer func() {
		if err := file.OutputFile.Close(); err != nil {
			c.app.Log.Error("file upload close failed", zap.Error(err))
		}
	}()
	for {
		req, err := stream.Recv()
		if file.FilePath == "" {
			file.SetFile("test.txt", "./")
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("file upload stream recv failed: %w", err)
		}
		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))
		c.app.Log.Debug("received a chunk with size", zap.Uint32("size", fileSize))
		if err := file.Write(chunk); err != nil {
			return fmt.Errorf("file upload write failed: %w", err)
		}
	}

	c.app.Log.Debug("file uploaded", zap.Uint32("size", fileSize))
	return stream.SendAndClose(&proto.FileUploadResponse{Size: &fileSize})
}

type File struct {
	FilePath   string
	OutputFile *os.File
}

func NewFile() *File {
	return &File{}
}

func (f *File) SetFile(fileName, path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("setFile mkdir failed: %w", err)
	}
	f.FilePath = filepath.Join(path, fileName)
	file, err := os.Create(f.FilePath)
	if err != nil {
		return fmt.Errorf("setFile create failed: %w", err)
	}
	f.OutputFile = file
	return nil
}

func (f *File) Write(chunk []byte) error {
	if f.OutputFile == nil {
		return nil
	}
	_, err := f.OutputFile.Write(chunk)
	return err
}

func (f *File) Close() error {
	return f.OutputFile.Close()
}
