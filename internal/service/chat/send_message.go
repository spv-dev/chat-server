package chat

import (
	"context"
	"fmt"

	"github.com/spv-dev/chat-server/internal/model"
	"github.com/spv-dev/chat-server/internal/validator"
)

// SendMessage проверяет сообщение и отправляет в слой БД
func (s *serv) SendMessage(ctx context.Context, info *model.MessageInfo) error {
	if info == nil {
		return fmt.Errorf("Пустая информация о сообщении")
	}
	if err := validator.CheckBody(info.Body); err != nil {
		return err
	}
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.chatRepository.SendMessage(ctx, info)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
