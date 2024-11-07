package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/spv-dev/chat-server/internal/api/chat"
	"github.com/spv-dev/chat-server/internal/model"
	"github.com/spv-dev/chat-server/internal/service"
	serviceMocks "github.com/spv-dev/chat-server/internal/service/mocks"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

func TestCreateChat(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.CreateChatRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		chatID  = gofakeit.Int64()
		userIDs = []int64{gofakeit.Int64(), gofakeit.Int64(), gofakeit.Int64()}
		title   = gofakeit.JobTitle()

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateChatRequest{
			Info: &desc.ChatInfo{
				Title:   title,
				UserIds: userIDs,
			},
		}

		resp = &desc.CreateChatResponse{
			Id: chatID,
		}

		chatInfo = &model.ChatInfo{
			Title:   title,
			UserIDs: userIDs,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateChatResponse
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "Success Create Chat",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: resp,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateChatMock.Expect(ctx, chatInfo).Return(chatID, nil)
				return mock
			},
		},
		{
			name: "Error Create Chat",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateChatMock.Expect(ctx, chatInfo).Return(0, serviceErr)
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

			res, err := api.CreateChat(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
