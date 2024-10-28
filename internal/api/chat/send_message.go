package chat

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/spv-dev/chat-server/internal/converter"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

// SendMessage отправка сообщения в чат
func (s *Server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	info := converter.ToMessageInfoFromDesc(req.GetInfo())
	err := s.chatService.SendMessage(ctx, &info)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
