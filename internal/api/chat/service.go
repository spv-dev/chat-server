package chat

import (
	"github.com/spv-dev/chat-server/internal/service"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

type Server struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService
}

func NewServer(chatService service.ChatService) *Server {
	return &Server{
		chatService: chatService,
	}
}
