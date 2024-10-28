package chat

import (
	"context"
	"fmt"
	"log"

	"github.com/spv-dev/chat-server/internal/model"
	"github.com/spv-dev/chat-server/internal/validator"
)

// CreateChat проверяет чат и отправляет на создание в слой БД
func (s *serv) CreateChat(ctx context.Context, info *model.ChatInfo) (int64, error) {
	if info == nil {
		return 0, fmt.Errorf("Пустая информация о чате")
	}
	if err := validator.CheckTitle(info.Title); err != nil {
		return 0, err
	}
	var id int64
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error
		log.Println("createchat service")
		id, errTx = s.chatRepository.CreateChat(ctx, info)
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
