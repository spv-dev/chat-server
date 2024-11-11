package chat

import (
	"github.com/spv-dev/platform_common/pkg/db"

	"github.com/spv-dev/chat-server/internal/repository"
)

const (
	tableName       = "chats"
	idColumn        = "id"
	titleColumn     = "title"
	stateColumn     = "state"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
	deletedAtColumn = "deleted_at"

	messagesTable = "messages"
	chatIDColumn  = "chat_id"
	userIDColumn  = "user_id"
	typeColumn    = "type_id"
	bodyColumn    = "body"

	chatUsersTable = "chat_users"
)

type repo struct {
	db db.Client
}

// NewRepository конструктор нового слоя repo
func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{
		db: db,
	}
}
