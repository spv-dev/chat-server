package chat

import (
	"github.com/spv-dev/chat-server/internal/service"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

// Server структура сервера
type Server struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService
}

// NewServer конструктор создания нового сервера
func NewServer(chatService service.ChatService) *Server {
	return &Server{
		chatService: chatService,
	}
}
