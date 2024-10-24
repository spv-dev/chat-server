package chat

import (
	"context"

	"github.com/spv-dev/chat-server/internal/converter"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateChat создаёт нового пользователя
func (s *Server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	_, err := s.chatService.SendMessage(ctx, converter.ToMessageInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
