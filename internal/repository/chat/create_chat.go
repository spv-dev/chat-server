package chat

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/spv-dev/chat-server/internal/client/db"
	"github.com/spv-dev/chat-server/internal/model"
)

// CreateChat создание чата в БД
func (r *repo) CreateChat(ctx context.Context, info *model.ChatInfo) (int64, error) {
	if info == nil {
		return 0, fmt.Errorf("Пустая структура при добавлении чата")
	}
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(titleColumn).
		Values(info.Title).
		Suffix("returning id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		QueryRaw: query,
		Name:     "chat_repository.Create",
	}

	var chatID int64
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID); err != nil {
		return 0, err
	}

	log.Printf("inserted new chat with id = %v", chatID)

	return chatID, nil
}
