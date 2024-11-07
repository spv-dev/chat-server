package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/spv-dev/platform_common/pkg/db"
	"github.com/stretchr/testify/require"

	dbMock "github.com/spv-dev/chat-server/internal/client/db/mocks"
	"github.com/spv-dev/chat-server/internal/model"
	"github.com/spv-dev/chat-server/internal/repository"
	repoMocks "github.com/spv-dev/chat-server/internal/repository/mocks"
	"github.com/spv-dev/chat-server/internal/service/chat"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.MessageInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		chatID = gofakeit.Int64()
		userID = gofakeit.Int64()
		text   = gofakeit.Name()

		req = &model.MessageInfo{
			Body:   text,
			ChatID: chatID,
			UserID: userID,
		}
		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name               string
		args               args
		err                error
		chatRepositoryMock chatRepositoryMockFunc
		dbMockFunc         txManagerMockFunc
	}{
		{
			name: "Success Send Message",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.SendMessageMock.Expect(ctx, req).Return(nil)
				return mock
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				mock := dbMock.NewTxManagerMock(t)
				mock.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})

				return mock
			},
		},
		{
			name: "Error Send Message",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: repoErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.SendMessageMock.Expect(ctx, req).Return(repoErr)
				return mock
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				mock := dbMock.NewTxManagerMock(t)
				mock.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			chatRepoMock := tt.chatRepositoryMock(mc)
			txManagerMock := tt.dbMockFunc(mc)
			service := chat.NewService(chatRepoMock, txManagerMock)

			err := service.SendMessage(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
