package step

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/arefev/gophkeeper/internal/client/app"
	mock_app "github.com/arefev/gophkeeper/internal/client/app/mock"
	"github.com/arefev/gophkeeper/internal/client/tui/model"
	"github.com/arefev/gophkeeper/internal/logger"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestStep(t *testing.T) {
	t.Run("start step success", func(t *testing.T) {
		var buf bytes.Buffer
		var in bytes.Buffer
		in.Write([]byte("q"))

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)
		app := app.NewApp(nil, l)

		p := tea.NewProgram(NewStart(app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
		require.Contains(t, buf.String(), "Добро пожаловать")
	})

	t.Run("login step success", func(t *testing.T) {
		var buf bytes.Buffer
		var in bytes.Buffer
		in.Write([]byte("q"))

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)
		app := app.NewApp(nil, l)

		p := tea.NewProgram(NewLogin(app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
		require.Contains(t, buf.String(), "## Авторизация")
	})

	t.Run("login action success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer
		in.Write([]byte("q"))

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		data := &model.LoginData{
			Login:    "test_login",
			Password: "test_password",
		}

		l, err := logger.Build("debug")
		require.NoError(t, err)

		conn := mock_app.NewMockConnection(ctrl)
		conn.EXPECT().Login(gomock.Any(), data.Login, data.Password).MinTimes(1).MaxTimes(1)
		conn.EXPECT().SetToken(gomock.Any()).MinTimes(1).MaxTimes(1)
		conn.EXPECT().CheckTokenCmd().MinTimes(1).MaxTimes(1)
		app := app.NewApp(conn, l)

		p := tea.NewProgram(NewLoginAction(data, app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
	})

	t.Run("reg step success", func(t *testing.T) {
		var buf bytes.Buffer
		var in bytes.Buffer
		in.Write([]byte("q"))

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)
		app := app.NewApp(nil, l)

		p := tea.NewProgram(NewReg(app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
		require.Contains(t, buf.String(), "## Регистрация")
	})

	t.Run("reg action success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer
		in.Write([]byte("q"))

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		data := &model.RegData{
			Login:    "test_login",
			Password: "test_password",
		}

		l, err := logger.Build("debug")
		require.NoError(t, err)

		conn := mock_app.NewMockConnection(ctrl)
		conn.EXPECT().Register(gomock.Any(), data.Login, data.Password).MinTimes(1).MaxTimes(1)
		conn.EXPECT().SetToken(gomock.Any()).MinTimes(1).MaxTimes(1)
		conn.EXPECT().CheckTokenCmd().MinTimes(1).MaxTimes(1)
		app := app.NewApp(conn, l)

		p := tea.NewProgram(NewRegAction(app, data), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
	})

	t.Run("lk step success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer
		in.Write([]byte("q"))

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)

		conn := mock_app.NewMockConnection(ctrl)
		conn.EXPECT().CheckTokenCmd().MinTimes(1).MaxTimes(1)
		app := app.NewApp(conn, l)

		p := tea.NewProgram(NewLK(app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
		require.Contains(t, buf.String(), "## Личный кабинет")
	})
}
