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

func TestCreateChat(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	chatID := gofakeit.Int64()
	text := gofakeit.Name()

	req := &model.ChatInfo{
		Title:   text,
		UserIDs: []int64{gofakeit.Int64(), gofakeit.Int64(), gofakeit.Int64()},
	}

	reqEmptyTitle := &model.ChatInfo{
		Title:   "",
		UserIDs: []int64{gofakeit.Int64(), gofakeit.Int64(), gofakeit.Int64()},
	}

	mc := minimock.NewController(t)

	repo := repoMocks.NewChatRepositoryMock(mc)
	trans := dbMock.NewTxManagerMock(mc)

	service := chat.NewService(repo, trans)

	t.Run("create chat success", func(t *testing.T) {
		t.Parallel()

		repo.CreateChatMock.Expect(ctx, req).Return(chatID, nil)
		repo.AddUsersToChatMock.Expect(ctx, chatID, req.UserIDs).Return(nil)
		trans.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
			return handler(ctx)
		})

		id, err := service.CreateChat(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, id, chatID)
	})

	repoErr := errors.New("repo error")
	t.Run("create chat error", func(t *testing.T) {
		t.Parallel()

		repo.CreateChatMock.Expect(ctx, req).Return(0, repoErr)
		trans.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
			return handler(ctx)
		})

		_, err := service.CreateChat(ctx, req)

		assert.Equal(t, err, repoErr)
	})

	errorAddUsers := errors.New("add users error")
	t.Run("create chat add users error", func(t *testing.T) {
		t.Parallel()

		repo.CreateChatMock.Expect(ctx, req).Return(chatID, nil)
		repo.AddUsersToChatMock.Expect(ctx, chatID, req.UserIDs).Return(errorAddUsers)
		trans.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
			return handler(ctx)
		})

		_, err := service.CreateChat(ctx, req)

		assert.Equal(t, err, errorAddUsers)
	})

	errorEmptyChatInfo := errors.New("Пустая информация о чате")
	t.Run("create chat empty info error", func(t *testing.T) {
		t.Parallel()

		_, err := service.CreateChat(ctx, nil)

		assert.Equal(t, err, errorEmptyChatInfo)
	})

	errorEmptyChatTitle := errors.New("Пустое наименование чата")
	t.Run("create chat empty title error", func(t *testing.T) {
		t.Parallel()

		_, err := service.CreateChat(ctx, reqEmptyTitle)

		assert.Equal(t, err, errorEmptyChatTitle)
	})
}
