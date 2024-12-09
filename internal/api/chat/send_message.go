package chat

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

// SendMessage отправка сообщения в чат
func (s *Server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	/*
		err := s.chatService.SendMessage(ctx, converter.ToMessageInfoFromDesc(req.GetInfo()))
		if err != nil {
			return nil, err
		}

		return nil, nil
	*/
	s.mxChannel.RLock()
	chatChan, ok := s.channels[req.Info.GetChatId()]
	s.mxChannel.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	chatChan <- req.GetInfo()

	return &emptypb.Empty{}, nil
}
