package chat

import (
	"github.com/spv-dev/chat-server/internal/client/db"
	"github.com/spv-dev/chat-server/internal/repository"
	"github.com/spv-dev/chat-server/internal/service"
)

type serv struct {
	chatRepository repository.ChatRepository
	txManager      db.TxManager
}

func NewService(chatRepository repository.ChatRepository,
	txManager db.TxManager) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
		txManager:      txManager,
	}
}
