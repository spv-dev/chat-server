package chat

import (
	"context"
	"fmt"

	"github.com/spv-dev/chat-server/internal/model"
	"github.com/spv-dev/chat-server/internal/validator"
)

// CreateChat проверяет чат и отправляет на создание в слой БД
func (s *serv) CreateChat(ctx context.Context, info *model.ChatInfo) (model.Chat, error) {
	if info == nil {
		return model.Chat{}, fmt.Errorf("Пустая информация о чате")
	}

	if err := validator.CheckTitle(info.Title); err != nil {
		return model.Chat{}, err
	}

	var chat model.Chat
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error
		chat, errTx = s.chatRepository.CreateChat(ctx, info)
		if errTx != nil {
			return errTx
		}

		if info.UserIDs != nil && len(info.UserIDs) > 0 {
			errTx = s.chatRepository.AddUsersToChat(ctx, chat.ID, info.UserIDs)
			if errTx != nil {
				return errTx
			}
		}

		return nil
	})

	if err != nil {
		return model.Chat{}, err
	}

	return chat, nil
}
