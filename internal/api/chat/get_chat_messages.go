package chat

import (
	"context"

	"github.com/spv-dev/chat-server/internal/converter"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

// GetChatMessages получение сообщений чата по идентификатору
func (s *Server) GetChatMessages(ctx context.Context, req *desc.GetChatMessagesRequest) (*desc.GetChatMessagesResponse, error) {
	messages, err := s.chatService.GetChatMessages(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &desc.GetChatMessagesResponse{
		Messages: converter.ToMessagesFromService(messages),
	}, nil
}
