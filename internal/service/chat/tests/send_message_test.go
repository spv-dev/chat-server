package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/spv-dev/platform_common/pkg/db"
	"github.com/stretchr/testify/assert"

	dbMock "github.com/spv-dev/chat-server/internal/client/db/mocks"
	"github.com/spv-dev/chat-server/internal/model"
	repoMocks "github.com/spv-dev/chat-server/internal/repository/mocks"
	"github.com/spv-dev/chat-server/internal/service/chat"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	chatID := gofakeit.Int64()
	userID := gofakeit.Int64()
	text := gofakeit.Name()

	req := &model.MessageInfo{
		Body:   text,
		ChatID: chatID,
		UserID: userID,
	}

	reqEmptyBody := &model.MessageInfo{
		Body:   "",
		ChatID: chatID,
		UserID: userID,
	}

	messObj := model.Message{
		ID:    gofakeit.Int64(),
		State: 1,
		Type:  10,
		Info: model.MessageInfo{
			ChatID: 10,
			UserID: 100500,
			Body:   "Test message",
		},
	}

	mc := minimock.NewController(t)

	repo := repoMocks.NewChatRepositoryMock(mc)
	trans := dbMock.NewTxManagerMock(mc)

	service := chat.NewService(repo, trans)

	t.Run("send message success", func(t *testing.T) {
		t.Parallel()

		repo.SendMessageMock.Expect(ctx, req).Return(messObj, nil)
		trans.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
			return handler(ctx)
		})

		message, err := service.SendMessage(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, messObj, message)
	})

	repoErr := errors.New("repo error")
	t.Run("send message error", func(t *testing.T) {
		t.Parallel()

		repo.SendMessageMock.Expect(ctx, req).Return(model.Message{}, repoErr)
		trans.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
			return handler(ctx)
		})

		_, err := service.SendMessage(ctx, req)

		assert.Equal(t, err, repoErr)
	})

	errEmptyMessageInfo := errors.New("Пустая информация о сообщении")
	t.Run("send message empty message error", func(t *testing.T) {
		t.Parallel()

		_, err := service.SendMessage(ctx, nil)

		assert.Equal(t, err, errEmptyMessageInfo)
	})

	errEmptyMessageBody := errors.New("Тело сообщения пустое")
	t.Run("send message empty message body error", func(t *testing.T) {
		t.Parallel()

		_, err := service.SendMessage(ctx, reqEmptyBody)

		assert.Equal(t, err, errEmptyMessageBody)
	})
}
