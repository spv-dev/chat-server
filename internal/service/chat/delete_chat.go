package chat

import (
	"context"
)

// DeleteChat удаление чата
func (s *serv) DeleteChat(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.chatRepository.DeleteChat(ctx, id)
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
