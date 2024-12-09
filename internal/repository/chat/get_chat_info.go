package chat

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/spv-dev/platform_common/pkg/db"

	"github.com/spv-dev/chat-server/internal/model"
)

// GetChatInfo получение информации о чате
func (r *repo) GetChatInfo(ctx context.Context, id int64) (model.Chat, error) {
	builder := sq.Select(idColumn, titleColumn, stateColumn, createdAtColumn).
		From(chatsTable).
		Where(sq.Eq{idColumn: id, stateColumn: 1}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return model.Chat{}, fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		QueryRaw: query,
		Name:     "chat_repository.GetChatInfo",
	}

	var chat model.Chat
	err = r.db.DB().ScanOneContext(ctx, &chat, q, args...)
	if err != nil {
		return model.Chat{}, fmt.Errorf("failed to get chat info: %v", err)
	}

	return chat, nil
}
