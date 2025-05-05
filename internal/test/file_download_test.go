package test

import (
	"context"
	"os"
	"testing"

	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/stretchr/testify/require"
)

func TestFileDownload(t *testing.T) {
	t.Run("file download success", func(t *testing.T) {
		const (
			login string = "test"
			pwd   string = "test"
			mName string = "from yandex"
			mType string = "creds"
			data  string = "login/password"
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

		err = conn.TextUpload(ctx, []byte(data), mName, mType)
		require.NoError(t, err)

		list, err := conn.GetList(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, list)

		path, err := conn.FileDownload(ctx, (*list)[0].UUID)
		require.NoError(t, err)

		app.Log.Sugar().Infof("file path %s", path)
		require.FileExists(t, path)

		b, err := os.ReadFile(path)
		require.NoError(t, err)
		require.Equal(t, b, []byte(data))

		err = os.Remove(path)
		require.NoError(t, err)
	})
}
