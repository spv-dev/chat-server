package chat

import (
	"sync"

	"github.com/spv-dev/chat-server/internal/service"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

type Chat struct {
	streams map[int64]desc.ChatV1_ConnectChatServer
	m       sync.RWMutex
}

// Server структура сервера
type Server struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService

	chats  map[int64]*Chat
	mxChat sync.RWMutex

	channels  map[int64]chan *desc.Message
	mxChannel sync.RWMutex
}

// NewServer конструктор создания нового сервера
func NewServer(chatService service.ChatService) *Server {
	return &Server{
		chatService: chatService,
		chats:       make(map[int64]*Chat),
		channels:    make(map[int64]chan *desc.Message),
	}
}
