package chat

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/spv-dev/platform_common/pkg/db"

	"github.com/spv-dev/chat-server/internal/model"
)

// SendMessage сохранение сообщения в БД
func (r *repo) SendMessage(ctx context.Context, info *model.MessageInfo) (model.Message, error) {
	if info == nil {
		return model.Message{}, fmt.Errorf("Пустая структура сообщения")
	}
	// добавим информацию о чате
	builder := sq.Insert(messagesTable).
		Columns(chatIDColumn, userIDColumn, bodyColumn).
		Values(info.ChatID, info.UserID, info.Body).
		Suffix("returning " + strings.Join([]string{idColumn, userIDColumn, bodyColumn, typeColumn, stateColumn, createdAtColumn, updatedAtColumn}, ", ")).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return model.Message{}, fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		QueryRaw: query,
		Name:     "chat_repository.SendMessage",
	}

	var message model.Message
	if err = r.db.DB().ScanOneContext(ctx, &message, q, args...); err != nil {
		return model.Message{}, err
	}
	return message, err
}
