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

func TestCreateChat(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	chatID := gofakeit.Int64()
	userIDs := []int64{gofakeit.Int64(), gofakeit.Int64(), gofakeit.Int64()}
	title := gofakeit.JobTitle()

	req := &desc.CreateChatRequest{
		Info: &desc.ChatInfo{
			Title:   title,
			UserIds: userIDs,
		},
	}

	resp := &desc.CreateChatResponse{
		Id: chatID,
	}

	chatInfo := &model.ChatInfo{
		Title:   title,
		UserIDs: userIDs,
	}

	mc := minimock.NewController(t)

	service := serviceMocks.NewChatServiceMock(mc)

	api := chat.NewServer(service)

	t.Run("create chat success", func(t *testing.T) {
		t.Parallel()

		service.CreateChatMock.Expect(ctx, chatInfo).Return(chatID, nil)

		r, err := api.CreateChat(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, r, resp)
	})

	serviceErr := errors.New("service error")
	t.Run("create chat error", func(t *testing.T) {
		t.Parallel()

		service.CreateChatMock.Expect(ctx, chatInfo).Return(0, serviceErr)

		_, err := api.CreateChat(ctx, req)

		assert.Equal(t, err, serviceErr)
	})
}
