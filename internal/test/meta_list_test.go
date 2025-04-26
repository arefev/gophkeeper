package test

import (
	"context"
	"testing"

	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/stretchr/testify/require"
)

func TestMetaList(t *testing.T) {
	t.Run("get list success", func(t *testing.T) {
		const (
			login string = "test"
			pwd   string = "test"
			mName string = "from yandex"
			mType string = "creds"
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

		err = conn.TextUpload(ctx, []byte("login/password"), mName, mType)
		require.NoError(t, err)

		list, err := conn.GetList(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, list)

		require.Equal(t, (*list)[0].Name, mName)
		require.Equal(t, (*list)[0].Type, mType)
	})

	t.Run("delete from list success", func(t *testing.T) {
		const (
			login string = "test"
			pwd   string = "test"
			mName string = "from yandex"
			mType string = "creds"
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

		err = conn.TextUpload(ctx, []byte("login/password"), mName, mType)
		require.NoError(t, err)

		list, err := conn.GetList(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, list)

		require.Equal(t, (*list)[0].Name, mName)
		require.Equal(t, (*list)[0].Type, mType)

		err = conn.Delete(ctx, (*list)[0].UUID)
		require.NoError(t, err)

		list, err = conn.GetList(ctx)
		require.NoError(t, err)
		require.Empty(t, list)
	})
}
