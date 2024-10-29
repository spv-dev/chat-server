package chat

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/spv-dev/chat-server/internal/client/db"
)

// AddUsersToChat добавление пользователей в чат
func (r *repo) AddUsersToChat(ctx context.Context, chatID int64, userIDs []int64) error {
	// добавим информацию о пользователях в чате
	builderUsers := sq.Insert(chatUsersTable).
		PlaceholderFormat(sq.Dollar).
		Columns(chatIDColumn, userIDColumn)

	for _, userID := range userIDs {
		builderUsers = builderUsers.Values(chatID, userID)
	}

	query, args, err := builderUsers.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert_users query: %v", err)
	}

	q := db.Query{
		QueryRaw: query,
		Name:     "chat_repository.AddUsersToChat",
	}

	if _, err = r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("failed to add users to chat: %v", err)
	}

	log.Printf("inserted new users to chat  = %v", chatID)

	return nil
}
