package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

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

	dt := time.Now()
	chatObj := desc.Chat{
		Id: chatID,
		Info: &desc.ChatInfo{
			Title: gofakeit.JobTitle(),
		},
		State:     1,
		CreatedAt: timestamppb.New(dt),
	}

	chatModel := model.Chat{
		ID: chatID,
		Info: model.ChatInfo{
			Title: gofakeit.JobTitle(),
		},
		State:     1,
		CreatedAt: dt,
	}

	resp := &desc.CreateChatResponse{
		Chat: &chatObj,
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

		service.CreateChatMock.Expect(ctx, chatInfo).Return(chatModel, nil)

		r, err := api.CreateChat(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, r, resp)
	})

	serviceErr := errors.New("service error")
	t.Run("create chat error", func(t *testing.T) {
		t.Parallel()

		service.CreateChatMock.Expect(ctx, chatInfo).Return(model.Chat{}, serviceErr)

		_, err := api.CreateChat(ctx, req)

		assert.Equal(t, err, serviceErr)
	})
}
