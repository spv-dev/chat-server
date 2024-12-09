package chat

import (
	"context"

	"github.com/spv-dev/chat-server/internal/model"
)

// GetChatInfo получает информацию о чате
func (s *serv) GetChatInfo(ctx context.Context, id int64) (model.Chat, error) {
	chat, err := s.chatRepository.GetChatInfo(ctx, id)
	if err != nil {
		return model.Chat{}, err
	}

	return chat, nil
}
