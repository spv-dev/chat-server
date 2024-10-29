package chat

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/spv-dev/chat-server/internal/client/db"
)

// DeleteChat удаление чата из БД
func (r *repo) DeleteChat(ctx context.Context, id int64) error {
	// будем не удалять информацию о чате, а менять статус
	builder := sq.Update(tableName).
		Set(stateColumn, 0).
		Set(deletedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		QueryRaw: query,
		Name:     "chat_repository.Delete",
	}

	if _, err = r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("failed to update chat: %v", err)
	}
	return nil
}
