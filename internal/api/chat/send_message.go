package chat

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/spv-dev/chat-server/internal/converter"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

// SendMessage отправка сообщения в чат
func (s *Server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*desc.SendMessageResponse, error) {

	message, err := s.chatService.SendMessage(ctx, converter.ToMessageInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	s.mxChannel.RLock()
	chatChan, ok := s.channels[req.Info.GetChatId()]
	s.mxChannel.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	msg := converter.ToMessageFromModel(&message)

	chatChan <- msg

	return &desc.SendMessageResponse{
		Message: msg,
	}, nil
}
