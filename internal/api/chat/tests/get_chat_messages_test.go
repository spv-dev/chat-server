package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/spv-dev/chat-server/internal/api/chat"
	"github.com/spv-dev/chat-server/internal/model"
	"github.com/spv-dev/chat-server/internal/service"
	serviceMocks "github.com/spv-dev/chat-server/internal/service/mocks"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

func TestGetChatMessages(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.GetChatMessagesRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		chatID   = gofakeit.Int64()
		userID   = gofakeit.Int64()
		state    = int32(1)
		messType = int32(10)

		id1 = gofakeit.Int64()
		id2 = gofakeit.Int64()
		id3 = gofakeit.Int64()

		text1 = gofakeit.Country()
		text2 = gofakeit.Country()
		text3 = gofakeit.Country()

		mt = time.Now()

		serviceErr = fmt.Errorf("service error")

		req = &desc.GetChatMessagesRequest{
			Id: chatID,
		}

		resp = &desc.GetChatMessagesResponse{
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

		messages = []*model.Message{
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
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetChatMessagesResponse
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "Success Get Chat Messages",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: resp,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.GetChatMessagesMock.Expect(ctx, chatID).Return(messages, nil)
				return mock
			},
		},
		{
			name: "Error Get Chat Messages",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.GetChatMessagesMock.Expect(ctx, chatID).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			chatServiceMock := tt.chatServiceMock(mc)
			api := chat.NewServer(chatServiceMock)

			res, err := api.GetChatMessages(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
