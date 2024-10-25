package chat

import (
	"context"

	"github.com/spv-dev/chat-server/internal/model"
)

// CreateChat проверяет чат и отправляет на создание в слой БД
func (s *serv) CreateChat(ctx context.Context, info *model.ChatInfo) (int64, error) {
	var id int64
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.chatRepository.CreateChat(ctx, info)
		if errTx != nil {
			return errTx
		}

		errTx = s.chatRepository.AddUsersToChat(ctx, id, info.UserIDs)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
