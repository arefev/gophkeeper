package service

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileService struct {
	Path     string
	MetaName string
	MetaType int
	Output   *os.File
}

func NewFile() *FileService {
	return &FileService{}
}

func (fs *FileService) SetFile(fileName, path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("setFile mkdir failed: %w", err)
	}
	fs.Path = filepath.Join(path, fileName)
	file, err := os.Create(fs.Path)
	if err != nil {
		return fmt.Errorf("setFile create failed: %w", err)
	}
	fs.Output = file
	return nil
}

func (fs *FileService) Write(chunk []byte) error {
	if fs.Output == nil {
		return nil
	}
	_, err := fs.Output.Write(chunk)
	return err
}

func (fs *FileService) Close() error {
	return fs.Output.Close()
}

func (fs *FileService) ReadAll() ([]byte, error) {
	b, err := os.ReadFile(fs.Path)
	if err != nil {
		return nil, fmt.Errorf("readall failed: %w", err)
	}
	return b, nil
}
