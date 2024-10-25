package chat

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/spv-dev/chat-server/internal/client/db"
	model "github.com/spv-dev/chat-server/internal/model"
	"github.com/spv-dev/chat-server/internal/repository"
	"github.com/spv-dev/chat-server/internal/repository/chat/converter"
	modelRepo "github.com/spv-dev/chat-server/internal/repository/chat/model"
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
)

type repo struct {
	db db.Client
}

// NewRepository конструктор слоя repo
func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{
		db: db,
	}
}

// CreateChat добавление чата в базу данных
func (r *repo) CreateChat(ctx context.Context, info *model.ChatInfo) (int64, error) {
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

// DeleteChat удаление чата из базы данных
// Не удаляет совсем, а устанавливает state = 0 /*DELETED*/
func (r *repo) DeleteChat(ctx context.Context, id int64) (*emptypb.Empty, error) {
	// будем не удалять информацию о чате, а менять статус
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(stateColumn, 0).
		Set(deletedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	q := db.Query{
		QueryRaw: query,
		Name:     "chat_repository.Delete",
	}

	if _, err = r.db.DB().ExecContext(ctx, q, args...); err != nil {
		log.Fatalf("failed to update chat: %v", err)
	}
	return nil, nil
}

// SendMessage добавляет сообщение в БД
func (r *repo) SendMessage(ctx context.Context, info *model.MessageInfo) (*emptypb.Empty, error) {
	// добавим информацию о чате
	builder := sq.Insert(messagesTable).
		PlaceholderFormat(sq.Dollar).
		Columns(chatIDColumn, userIDColumn, bodyColumn).
		Values(info.ChatID, info.UserID, info.Body).
		Suffix("returning id")

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	q := db.Query{
		QueryRaw: query,
		Name:     "chat_repository.SendMessage",
	}

	if _, err = r.db.DB().ExecContext(ctx, q, args...); err != nil {
		log.Fatalf("failed to send message: %v", err)
	}
	return nil, nil
}

// GetChatMessages получает сообщения указанного чата из БД
func (r *repo) GetChatMessages(ctx context.Context, id int64) ([]*model.Message, error) {
	// будем не удалять информацию о чате, а менять статус
	builder := sq.Select(idColumn, userIDColumn, bodyColumn, createdAtColumn, typeColumn).
		PlaceholderFormat(sq.Dollar).
		From(messagesTable).
		Where(sq.Eq{chatIDColumn: id, stateColumn: 1}).
		OrderBy("created_at DESC").
		Limit(20)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	q := db.Query{
		QueryRaw: query,
		Name:     "chat_repository.GetChatMessages",
	}

	var messagesRepo []*modelRepo.Message
	err = r.db.DB().ScanAllContext(ctx, &messagesRepo, q, args...)
	if err != nil {
		log.Fatalf("failed to get messages: %v", err)
	}

	return converter.ToMessagesFromRepo(messagesRepo), nil
}
