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

	mc := minimock.NewController(t)

	repo := repoMocks.NewChatRepositoryMock(mc)
	trans := dbMock.NewTxManagerMock(mc)

	service := chat.NewService(repo, trans)

	t.Run("send message success", func(t *testing.T) {
		t.Parallel()

		repo.SendMessageMock.Expect(ctx, req).Return(nil)
		trans.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
			return handler(ctx)
		})

		err := service.SendMessage(ctx, req)

		assert.NoError(t, err)
	})

	repoErr := errors.New("repo error")
	t.Run("send message error", func(t *testing.T) {
		t.Parallel()

		repo.SendMessageMock.Expect(ctx, req).Return(repoErr)
		trans.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
			return handler(ctx)
		})

		err := service.SendMessage(ctx, req)

		assert.Equal(t, err, repoErr)
	})

	errEmptyMessageInfo := errors.New("Пустая информация о сообщении")
	t.Run("send message empty message error", func(t *testing.T) {
		t.Parallel()

		err := service.SendMessage(ctx, nil)

		assert.Equal(t, err, errEmptyMessageInfo)
	})

	errEmptyMessageBody := errors.New("Тело сообщения пустое")
	t.Run("send message empty message body error", func(t *testing.T) {
		t.Parallel()

		err := service.SendMessage(ctx, reqEmptyBody)

		assert.Equal(t, err, errEmptyMessageBody)
	})
}
