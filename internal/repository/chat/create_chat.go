package chat

import (
	"context"
	"fmt"
	"log"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/spv-dev/platform_common/pkg/db"

	"github.com/spv-dev/chat-server/internal/model"
)

// CreateChat создание чата в БД
func (r *repo) CreateChat(ctx context.Context, info *model.ChatInfo) (model.Chat, error) {
	if info == nil {
		return model.Chat{}, fmt.Errorf("Пустая структура при добавлении чата")
	}
	builder := sq.Insert(chatsTable).
		Columns(titleColumn).
		Values(info.Title).
		Suffix("returning " + strings.Join([]string{idColumn, titleColumn, stateColumn, createdAtColumn}, ", ")).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return model.Chat{}, err
	}

	q := db.Query{
		QueryRaw: query,
		Name:     "chat_repository.Create",
	}

	var chat model.Chat
	//if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chat); err != nil {
	if err = r.db.DB().ScanOneContext(ctx, &chat, q, args...); err != nil {
		return model.Chat{}, err
	}

	log.Printf("inserted new chat with id = %v", chat.ID)

	return chat, nil
}
