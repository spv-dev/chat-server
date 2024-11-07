package chat

import (
	"github.com/spv-dev/platform_common/pkg/db"

	"github.com/spv-dev/chat-server/internal/repository"
	"github.com/spv-dev/chat-server/internal/service"
)

type serv struct {
	chatRepository repository.ChatRepository
	txManager      db.TxManager
}

// NewService конструктор нового сервиса
func NewService(chatRepository repository.ChatRepository,
	txManager db.TxManager) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
		txManager:      txManager,
	}
}
