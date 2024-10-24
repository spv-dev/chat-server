package chat

import (
	"context"

	"github.com/spv-dev/chat-server/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateChat проверяет чат и отправляет на создание в слой БД
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
