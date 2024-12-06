package chat

import (
	"context"

	"github.com/spv-dev/chat-server/internal/converter"
	"github.com/spv-dev/chat-server/internal/logger"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
	"go.uber.org/zap"
)

// CreateChat создаёт нового пользователя
func (s *Server) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	logger.Info("Create Chat...", zap.String("Title", req.Info.GetTitle()))
	logger.Debug("Create Chat...", zap.String("Title", req.Info.GetTitle()))
	id, err := s.chatService.CreateChat(ctx, converter.ToChatInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return &desc.CreateChatResponse{
		Id: id,
	}, nil
}
