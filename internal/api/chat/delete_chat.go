package chat

import (
	"context"

	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateChat создаёт нового пользователя
func (s *Server) DeleteChat(ctx context.Context, req *desc.DeleteChatRequest) (*emptypb.Empty, error) {
	_, err := s.chatService.DeleteChat(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
