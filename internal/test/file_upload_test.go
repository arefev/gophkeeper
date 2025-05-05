package test

import (
	"context"
	"os"
	"testing"

	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/stretchr/testify/require"
)

func TestFileUpload(t *testing.T) {
	t.Run("file upload success", func(t *testing.T) {
		const (
			login string = "test"
			pwd   string = "test"
		)
		ctx := context.Background()
		p := NewPrepare()
		defer func() {
			require.NoError(t, p.Close(ctx))
		}()

		require.NoError(t, p.runDB(ctx))

		app, err := p.app()
		require.NoError(t, err)

		require.NoError(t, p.server(app))

		conn := connection.NewGRPCClient(app.Conf.ChunkSize, app.Log)
		require.NoError(t, conn.Connect(app.Conf.Address))

		defer func() {
			require.NoError(t, conn.Close())
		}()

		token, err := conn.Register(ctx, login, pwd)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		conn.SetToken(token)

		dir, err := os.Getwd()
		require.NoError(t, err)

		file := dir + "/files/test.txt"

		err = conn.FileUpload(ctx, file, "test", "test")
		require.NoError(t, err)
	})
}
