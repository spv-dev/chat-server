package chat

import (
	"context"

	"github.com/spv-dev/chat-server/internal/model"
)

func (s *serv) GetChatMessages(ctx context.Context, id int64) ([]*model.Message, error) {
	messages, err := s.chatRepository.GetChatMessages(ctx, id)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
