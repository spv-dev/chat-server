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

func TestGetChatMessages(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	chatID := gofakeit.Int64()
	userID := gofakeit.Int64()
	state := int32(1)
	messType := int32(10)

	id1 := gofakeit.Int64()
	id2 := gofakeit.Int64()
	id3 := gofakeit.Int64()

	text1 := gofakeit.Country()
	text2 := gofakeit.Country()
	text3 := gofakeit.Country()

	mt := time.Now()

	req := &desc.GetChatMessagesRequest{
		Id: chatID,
	}

	resp := &desc.GetChatMessagesResponse{
		Messages: []*desc.Message{
			{
				Id:        id1,
				State:     state,
				Type:      messType,
				CreatedAt: timestamppb.New(mt),
				Info: &desc.MessageInfo{
					ChatId: chatID,
					UserId: userID,
					Body:   text1,
				},
			},
			{
				Id:        id2,
				State:     state,
				Type:      messType,
				CreatedAt: timestamppb.New(mt),
				Info: &desc.MessageInfo{
					ChatId: chatID,
					UserId: userID,
					Body:   text2,
				},
			},
			{
				Id:        id3,
				State:     state,
				Type:      messType,
				CreatedAt: timestamppb.New(mt),
				Info: &desc.MessageInfo{
					ChatId: chatID,
					UserId: userID,
					Body:   text3,
				},
			},
		},
	}

	messages := []*model.Message{
		{
			ID:        id1,
			State:     state,
			Type:      messType,
			CreatedAt: mt,
			Info: model.MessageInfo{
				ChatID: chatID,
				UserID: userID,
				Body:   text1,
			},
		},
		{
			ID:        id2,
			State:     state,
			Type:      messType,
			CreatedAt: mt,
			Info: model.MessageInfo{
				ChatID: chatID,
				UserID: userID,
				Body:   text2,
			},
		},
		{
			ID:        id3,
			State:     state,
			Type:      messType,
			CreatedAt: mt,
			Info: model.MessageInfo{
				ChatID: chatID,
				UserID: userID,
				Body:   text3,
			},
		},
	}

	mc := minimock.NewController(t)

	service := serviceMocks.NewChatServiceMock(mc)

	api := chat.NewServer(service)

	t.Run("get chat messages success", func(t *testing.T) {
		t.Parallel()

		service.GetChatMessagesMock.Expect(ctx, chatID).Return(messages, nil)

		res, err := api.GetChatMessages(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, res, resp)
	})

	serviceErr := errors.New("service error")
	t.Run("get chat messages error", func(t *testing.T) {
		t.Parallel()

		service.GetChatMessagesMock.Expect(ctx, chatID).Return(nil, serviceErr)

		_, err := api.GetChatMessages(ctx, req)

		assert.Equal(t, err, serviceErr)
	})
}
