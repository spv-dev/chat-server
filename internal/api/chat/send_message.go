package chat

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/spv-dev/chat-server/internal/converter"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

// SendMessage отправка сообщения в чат
func (s *Server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	err := s.chatService.SendMessage(ctx, converter.ToMessageInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
