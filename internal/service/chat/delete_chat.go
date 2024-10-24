package chat

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serv) DeleteChat(ctx context.Context, id int64) (*emptypb.Empty, error) {
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error

		_, errTx = s.chatRepository.DeleteChat(ctx, id)
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
