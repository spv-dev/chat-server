package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	dbMock "github.com/spv-dev/chat-server/internal/client/db/mocks"
	repoMocks "github.com/spv-dev/chat-server/internal/repository/mocks"
	"github.com/spv-dev/chat-server/internal/service/chat"
	"github.com/spv-dev/platform_common/pkg/db"
)

func TestDeleteChat(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	chatID := gofakeit.Int64()

	req := chatID

	mc := minimock.NewController(t)

	repo := repoMocks.NewChatRepositoryMock(mc)
	trans := dbMock.NewTxManagerMock(mc)

	service := chat.NewService(repo, trans)

	t.Run("delete chat success", func(t *testing.T) {
		t.Parallel()

		repo.DeleteChatMock.Expect(ctx, chatID).Return(nil)
		trans.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
			return handler(ctx)
		})

		err := service.DeleteChat(ctx, req)

		assert.NoError(t, err)
	})

	repoErr := errors.New("repo error")
	t.Run("delete chat error", func(t *testing.T) {
		t.Parallel()

		repo.DeleteChatMock.Expect(ctx, chatID).Return(repoErr)
		trans.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
			return handler(ctx)
		})

		err := service.DeleteChat(ctx, req)

		assert.Equal(t, err, repoErr)
	})
}
