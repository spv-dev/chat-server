package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/spv-dev/chat-server/internal/api/chat"
	serviceMocks "github.com/spv-dev/chat-server/internal/service/mocks"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

func TestDeleteChat(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	chatID := gofakeit.Int64()

	req := &desc.DeleteChatRequest{
		Id: chatID,
	}

	mc := minimock.NewController(t)

	service := serviceMocks.NewChatServiceMock(mc)

	api := chat.NewServer(service)

	t.Run("delete chat success", func(t *testing.T) {
		t.Parallel()

		service.DeleteChatMock.Expect(ctx, chatID).Return(nil)

		_, err := api.DeleteChat(ctx, req)

		assert.NoError(t, err)
	})

	serviceErr := errors.New("service error")
	t.Run("delete chat error", func(t *testing.T) {
		t.Parallel()

		service.DeleteChatMock.Expect(ctx, chatID).Return(serviceErr)

		_, err := api.DeleteChat(ctx, req)

		assert.Equal(t, err, serviceErr)
	})
}
