package chat

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/spv-dev/chat-server/internal/model"
)

// SendMessage проверяет сообщение и отправляет на создание в слой БД
func (s *serv) SendMessage(ctx context.Context, info *model.MessageInfo) (*emptypb.Empty, error) {
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error
		_, errTx = s.chatRepository.SendMessage(ctx, info)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}
