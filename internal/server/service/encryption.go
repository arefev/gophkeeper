package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"github.com/arefev/gophkeeper/internal/server/application"
)

type encryptionService struct {
	app *application.App
}

func NewEncryptionService(app *application.App) *encryptionService {
	return &encryptionService{
		app: app,
	}
}

func (enc *encryptionService) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(enc.app.Conf.EncryptionSecret))
	if err != nil {
		return nil, fmt.Errorf("ecrypt: new cipher failed: %w", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("ecrypt: read failed: %w", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

func (enc *encryptionService) Decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(enc.app.Conf.EncryptionSecret))
	if err != nil {
		return nil, fmt.Errorf("decrypt: new cipher failed: %w", err)
	}

	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return data, nil
}
