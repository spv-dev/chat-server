package chat

import (
	"context"

	"github.com/spv-dev/chat-server/internal/converter"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

// CreateChat создаёт нового пользователя
func (s *Server) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	info := converter.ToChatInfoFromDesc(req.GetInfo())
	id, err := s.chatService.CreateChat(ctx, &info)
	if err != nil {
		return nil, err
	}

	return &desc.CreateChatResponse{
		Id: id,
	}, nil
}
