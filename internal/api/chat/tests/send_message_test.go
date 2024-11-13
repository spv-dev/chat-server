package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/spv-dev/chat-server/internal/api/chat"
	"github.com/spv-dev/chat-server/internal/model"
	serviceMocks "github.com/spv-dev/chat-server/internal/service/mocks"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	chatID := gofakeit.Int64()
	userID := gofakeit.Int64()
	text := gofakeit.Company()

	req := &desc.SendMessageRequest{
		Info: &desc.MessageInfo{
			Body:   text,
			ChatId: chatID,
			UserId: userID,
		},
	}

	message := &model.MessageInfo{
		Body:   text,
		ChatID: chatID,
		UserID: userID,
	}

	mc := minimock.NewController(t)

	service := serviceMocks.NewChatServiceMock(mc)

	api := chat.NewServer(service)

	t.Run("send message success", func(t *testing.T) {
		t.Parallel()

		service.SendMessageMock.Expect(ctx, message).Return(nil)

		_, err := api.SendMessage(ctx, req)

		assert.NoError(t, err)
	})

	serviceErr := errors.New("service error")
	t.Run("send message error", func(t *testing.T) {
		t.Parallel()

		service.SendMessageMock.Expect(ctx, message).Return(serviceErr)

		_, err := api.SendMessage(ctx, req)

		assert.Equal(t, err, serviceErr)
	})
}
