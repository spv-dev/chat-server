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

func TestCreateChat(t *testing.T) {
	//t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.ChatInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		text = gofakeit.Name()

		req = &model.ChatInfo{
			Title:   text,
			UserIDs: []int64{gofakeit.Int64(), gofakeit.Int64(), gofakeit.Int64()},
		}

		resp          = gofakeit.Int64()
		repoErr       = fmt.Errorf("repo error")
		emptyErr      = fmt.Errorf("Пустая информация о чате")
		emptyTitleErr = fmt.Errorf("Пустое наименование чата")
		addUsersError = fmt.Errorf("error add user")
	)

	tests := []struct {
		name               string
		args               args
		wait               int64
		err                error
		chatRepositoryMock chatRepositoryMockFunc
		dbMockFunc         txManagerMockFunc
	}{
		{
			name: "Success Create Chat",
			args: args{
				ctx: ctx,
				req: req,
			},
			wait: resp,
			err:  nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.AddUsersToChatMock.Expect(ctx, resp, req.UserIDs).Return(nil)
				mock.CreateChatMock.Expect(ctx, req).Return(resp, nil)
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
			name: "Error Add Users",
			args: args{
				ctx: ctx,
				req: req,
			},
			wait: 0,
			err:  addUsersError,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.CreateChatMock.Expect(ctx, req).Return(resp, nil)
				mock.AddUsersToChatMock.Expect(ctx, resp, req.UserIDs).Return(addUsersError)

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
			name: "Empty Chat Info",
			args: args{
				ctx: ctx,
				req: nil,
			},
			wait: 0,
			err:  emptyErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				return repoMocks.NewChatRepositoryMock(mc)
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
		},
		{
			name: "Empty Title",
			args: args{
				ctx: ctx,
				req: &model.ChatInfo{
					Title:   "",
					UserIDs: []int64{23, 33, 44},
				},
			},
			wait: 0,
			err:  emptyTitleErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				return repoMocks.NewChatRepositoryMock(mc)
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
		},
		{
			name: "Error Create Chat",
			args: args{
				ctx: ctx,
				req: req,
			},
			wait: 0,
			err:  repoErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.CreateChatMock.Expect(ctx, req).Return(0, repoErr)
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
			//t.Parallel()
			chatRepoMock := tt.chatRepositoryMock(mc)
			txManagerMock := tt.dbMockFunc(mc)
			service := chat.NewService(chatRepoMock, txManagerMock)

			id, err := service.CreateChat(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.wait, id)
		})
	}
}
