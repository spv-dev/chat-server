package chat

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/spv-dev/platform_common/pkg/db"

	"github.com/spv-dev/chat-server/internal/model"
)

// GetChatMessages получение сообщений чата из БД по идентификатору чата
func (r *repo) GetChatMessages(ctx context.Context, id int64) ([]*model.Message, error) {
	// будем не удалять информацию о чате, а менять статус
	builder := sq.Select(idColumn, userIDColumn, bodyColumn, createdAtColumn, typeColumn, stateColumn).
		From(messagesTable).
		Where(sq.Eq{chatIDColumn: id, stateColumn: 1}).
		OrderBy("created_at DESC").
		Limit(20).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return []*model.Message{}, fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		QueryRaw: query,
		Name:     "chat_repository.SendMessage",
	}

	var messages []*model.Message
	err = r.db.DB().ScanAllContext(ctx, &messages, q, args...)
	if err != nil {
		return []*model.Message{}, fmt.Errorf("failed to get chat messages: %v", err)
	}

	return messages, nil
}
