package chat

import (
	"context"
	"fmt"

	"github.com/spv-dev/chat-server/internal/model"
	"github.com/spv-dev/chat-server/internal/validator"
)

// SendMessage проверяет сообщение и отправляет в слой БД
func (s *serv) SendMessage(ctx context.Context, info *model.MessageInfo) (model.Message, error) {
	if info == nil {
		return model.Message{}, fmt.Errorf("Пустая информация о сообщении")
	}

	if err := validator.CheckBody(info.Body); err != nil {
		return model.Message{}, err
	}

	var message model.Message
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error
		message, errTx = s.chatRepository.SendMessage(ctx, info)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return model.Message{}, err
	}

	return message, nil
}
