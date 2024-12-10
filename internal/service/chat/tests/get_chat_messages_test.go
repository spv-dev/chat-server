package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	dbMock "github.com/spv-dev/chat-server/internal/client/db/mocks"
	"github.com/spv-dev/chat-server/internal/model"
	repoMocks "github.com/spv-dev/chat-server/internal/repository/mocks"
	"github.com/spv-dev/chat-server/internal/service/chat"
)

func TestGetChatMessages(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	chatID := gofakeit.Int64()
	userID := gofakeit.Int64()
	state := int32(1)
	messType := int32(10)
	limit := gofakeit.Uint64()
	offset := gofakeit.Uint64()

	id1 := gofakeit.Int64()
	id2 := gofakeit.Int64()
	id3 := gofakeit.Int64()

	text1 := gofakeit.Country()
	text2 := gofakeit.Country()
	text3 := gofakeit.Country()

	mt := time.Now()

	req := chatID

	resp := []*model.Message{
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

	repo := repoMocks.NewChatRepositoryMock(mc)
	trans := dbMock.NewTxManagerMock(mc)

	service := chat.NewService(repo, trans)

	t.Run("get chat messages success", func(t *testing.T) {
		t.Parallel()

		repo.GetChatMessagesMock.Expect(ctx, chatID, limit, offset).Return(resp, nil)

		res, err := service.GetChatMessages(ctx, req, limit, offset)

		assert.NoError(t, err)
		assert.Equal(t, res, resp)
	})

	repoErr := errors.New("repo error")
	t.Run("get chat messages error", func(t *testing.T) {
		t.Parallel()

		repo.GetChatMessagesMock.Expect(ctx, chatID, limit, offset).Return(nil, repoErr)

		_, err := service.GetChatMessages(ctx, req, limit, offset)

		assert.Equal(t, err, repoErr)
	})

}
