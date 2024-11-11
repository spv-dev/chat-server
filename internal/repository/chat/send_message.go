package chat

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/spv-dev/platform_common/pkg/db"

	"github.com/spv-dev/chat-server/internal/model"
)

// SendMessage сохранение сообщения в БД
func (r *repo) SendMessage(ctx context.Context, info *model.MessageInfo) error {
	if info == nil {
		return fmt.Errorf("Пустая структура сообщения")
	}
	// добавим информацию о чате
	builder := sq.Insert(messagesTable).
		Columns(chatIDColumn, userIDColumn, bodyColumn).
		Values(info.ChatID, info.UserID, info.Body).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		QueryRaw: query,
		Name:     "chat_repository.SendMessage",
	}

	if _, err = r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}
