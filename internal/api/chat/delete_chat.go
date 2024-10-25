package chat

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

// DeleteChat удаление чата
func (s *Server) DeleteChat(ctx context.Context, req *desc.DeleteChatRequest) (*emptypb.Empty, error) {
	_, err := s.chatService.DeleteChat(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
