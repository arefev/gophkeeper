package step

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	application "github.com/arefev/gophkeeper/internal/client/app"
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

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)
		app := application.NewApp(nil, l)

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

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)
		app := application.NewApp(nil, l)

		p := tea.NewProgram(NewLogin(app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
		require.Contains(t, buf.String(), "Авторизация")
	})

	t.Run("login action success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer

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
		app := application.NewApp(conn, l)

		p := tea.NewProgram(NewLoginAction(data, app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
	})

	t.Run("reg step success", func(t *testing.T) {
		var buf bytes.Buffer
		var in bytes.Buffer

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)
		app := application.NewApp(nil, l)

		p := tea.NewProgram(NewReg(app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
		require.Contains(t, buf.String(), "Регистрация")
	})

	t.Run("reg action success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer

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
		app := application.NewApp(conn, l)

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

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)

		conn := mock_app.NewMockConnection(ctrl)
		conn.EXPECT().CheckTokenCmd().MinTimes(1).MaxTimes(1)
		app := application.NewApp(conn, l)

		p := tea.NewProgram(NewLK(app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
		require.Contains(t, buf.String(), "Личный кабинет")
	})

	t.Run("lk types step success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)

		conn := mock_app.NewMockConnection(ctrl)
		conn.EXPECT().CheckTokenCmd().MinTimes(1).MaxTimes(1)
		app := application.NewApp(conn, l)

		p := tea.NewProgram(NewLKTypes(app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
		require.Contains(t, buf.String(), "Какие данные вы хотите отправить?")
	})

	t.Run("lk form creds step success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)

		conn := mock_app.NewMockConnection(ctrl)
		conn.EXPECT().CheckTokenCmd().MinTimes(1).MaxTimes(1)
		app := application.NewApp(conn, l)

		p := tea.NewProgram(NewLKFormCreds(app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
		require.Contains(t, buf.String(), "Введите логин/пароль для сохранения")
	})

	t.Run("lk form bank step success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)

		conn := mock_app.NewMockConnection(ctrl)
		conn.EXPECT().CheckTokenCmd().MinTimes(1).MaxTimes(1)
		app := application.NewApp(conn, l)

		p := tea.NewProgram(NewLKFormBank(app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
		require.Contains(t, buf.String(), "Введите данные банковской карты для сохранения")
	})

	t.Run("lk form file step success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)

		conn := mock_app.NewMockConnection(ctrl)
		conn.EXPECT().CheckTokenCmd().MinTimes(1).MaxTimes(1)
		app := application.NewApp(conn, l)

		p := tea.NewProgram(NewLKFormFile(app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
		require.Contains(t, buf.String(), "Введите путь до файла для сохранения")
	})

	t.Run("creds send action step success", func(t *testing.T) {
		const mType string = "creds"
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer

		m := &model.CredsData{
			Name:     "name_test",
			Login:    "login_test",
			Password: "password_test",
		}
		data := fmt.Sprintf("Login: %s\nPassword: %s", m.Login, m.Password)

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)

		conn := mock_app.NewMockConnection(ctrl)
		conn.EXPECT().CheckTokenCmd().MinTimes(1).MaxTimes(2)
		conn.EXPECT().TextUpload(gomock.Any(), []byte(data), m.Name, mType).MinTimes(1).MaxTimes(1)
		app := application.NewApp(conn, l)

		p := tea.NewProgram(NewCredsSendAction(m, app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
	})

	t.Run("bank send action step success", func(t *testing.T) {
		const mType string = "card"
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer

		m := &model.BankData{
			Name:   "name_test",
			Number: "number_test",
			Exp:    "exp_test",
			CVV:    "cvv_test",
		}
		data := fmt.Sprintf("Number: %s\nExpired: %s\nCVV: %s", m.Number, m.Exp, m.CVV)

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)

		conn := mock_app.NewMockConnection(ctrl)
		conn.EXPECT().CheckTokenCmd().MinTimes(1).MaxTimes(2)
		conn.EXPECT().TextUpload(gomock.Any(), []byte(data), m.Name, mType).MinTimes(1).MaxTimes(1)
		app := application.NewApp(conn, l)

		p := tea.NewProgram(NewBankSendAction(m, app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
	})

	t.Run("file send action step success", func(t *testing.T) {
		const mType string = "file"
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer

		m := &model.FileData{
			Name: "name_test",
			Path: "path_test",
		}

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)

		conn := mock_app.NewMockConnection(ctrl)
		conn.EXPECT().CheckTokenCmd().MinTimes(1).MaxTimes(2)
		conn.EXPECT().FileUpload(gomock.Any(), m.Path, m.Name, mType).MinTimes(1).MaxTimes(1)
		app := application.NewApp(conn, l)

		p := tea.NewProgram(NewFileSendAction(m, app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
	})

	t.Run("lk list step success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer

		data := &[]model.MetaListData{
			{
				UUID: "uuid_test",
				Type: "type_test",
				Name: "name_test",
				File: "file_test",
				Date: "date_test",
			},
		}

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		l, err := logger.Build("debug")
		require.NoError(t, err)

		conn := mock_app.NewMockConnection(ctrl)
		conn.EXPECT().CheckTokenCmd().MinTimes(1).MaxTimes(2)
		conn.EXPECT().GetList(gomock.Any()).MinTimes(1).MaxTimes(1).Return(data, nil)
		app := application.NewApp(conn, l)

		p := tea.NewProgram(NewLKList(app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
	})

	t.Run("delete action step success", func(t *testing.T) {
		const uuid string = "uuid_test"
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		data := &[]model.MetaListData{
			{
				UUID: "uuid_test",
				Type: "type_test",
				Name: "name_test",
				File: "file_test",
				Date: "date_test",
			},
		}

		l, err := logger.Build("debug")
		require.NoError(t, err)

		conn := mock_app.NewMockConnection(ctrl)
		conn.EXPECT().CheckTokenCmd().MinTimes(1).MaxTimes(2)
		conn.EXPECT().Delete(gomock.Any(), uuid).MinTimes(1).MaxTimes(1).Return(nil)
		conn.EXPECT().GetList(gomock.Any()).MinTimes(1).MaxTimes(1).Return(data, nil)
		app := application.NewApp(conn, l)

		p := tea.NewProgram(NewDeleteAction(uuid, app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
	})

	t.Run("download action step success", func(t *testing.T) {
		const (
			uuid string = "uuid_test"
			path string = "path_test"
		)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buf bytes.Buffer
		var in bytes.Buffer

		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancel()

		data := &[]model.MetaListData{
			{
				UUID: "uuid_test",
				Type: "type_test",
				Name: "name_test",
				File: "file_test",
				Date: "date_test",
			},
		}

		l, err := logger.Build("debug")
		require.NoError(t, err)

		conn := mock_app.NewMockConnection(ctrl)
		conn.EXPECT().CheckTokenCmd().MinTimes(1).MaxTimes(2)
		conn.EXPECT().FileDownload(gomock.Any(), uuid).MinTimes(1).MaxTimes(1).Return(path, nil)
		conn.EXPECT().GetList(gomock.Any()).MinTimes(1).MaxTimes(1).Return(data, nil)
		app := application.NewApp(conn, l)

		p := tea.NewProgram(NewDownloadAction(uuid, app), tea.WithInput(&in), tea.WithOutput(&buf), tea.WithContext(ctx))
		if _, err := p.Run(); err != nil {
			require.Contains(t, err.Error(), "context deadline exceeded")
		}

		require.NotEmpty(t, buf.Len())
	})
}
