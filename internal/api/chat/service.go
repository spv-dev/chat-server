package chat

import (
	"github.com/spv-dev/chat-server/internal/service"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

// Server структрура сервиса
type Server struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService
}

// NewServer конструктор нового сервиса
func NewServer(chatService service.ChatService) *Server {
	return &Server{
		chatService: chatService,
	}
}
