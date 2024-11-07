package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

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

func TestGetChatMessages(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req int64
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

		req = chatID

		resp = []*model.Message{
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
		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name               string
		args               args
		wait               []*model.Message
		err                error
		chatRepositoryMock chatRepositoryMockFunc
		dbMockFunc         txManagerMockFunc
	}{
		{
			name: "Success Get Chat Messages",
			args: args{
				ctx: ctx,
				req: req,
			},
			wait: resp,
			err:  nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.GetChatMessagesMock.Expect(ctx, req).Return(resp, nil)
				return mock
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
		},
		{
			name: "Error Get Chat Messages",
			args: args{
				ctx: ctx,
				req: req,
			},
			wait: nil,
			err:  repoErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.GetChatMessagesMock.Expect(ctx, req).Return(nil, repoErr)
				return mock
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
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

			id, err := service.GetChatMessages(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.wait, id)
		})
	}
}
