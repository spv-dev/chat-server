package chat

import (
	"context"

	"go.uber.org/zap"

	"github.com/spv-dev/chat-server/internal/converter"
	"github.com/spv-dev/chat-server/internal/logger"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

// CreateChat создаёт новый чат
func (s *Server) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	logger.Info("Create Chat...", zap.String("Title", req.Info.GetTitle()))
	logger.Debug("Create Chat...", zap.String("Title", req.Info.GetTitle()))
	chat, err := s.chatService.CreateChat(ctx, converter.ToChatInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	s.channels[chat.ID] = make(chan *desc.Message, 100)

	return &desc.CreateChatResponse{
		Chat: converter.ToChatFromModel(&chat),
	}, nil
}
