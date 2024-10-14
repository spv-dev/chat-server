package database

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

const (
	dbDSN = "host=localhost port=54321 dbname=chat user=chat-user password=chatpwd sslmode=disable"
)

var pool *pgxpool.Pool

// InitDB инициализация базы данных
func InitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	pool = p
}

// CloseDB закрытие базы данных
func CloseDB() {
	pool.Close()
}

// CreateChatDB добавляет информацию о чате в Базу данных
func CreateChatDB(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	// добавим информацию о чате
	builder := sq.Insert("chats").
		PlaceholderFormat(sq.Dollar).
		Columns("title").
		Values(req.Info.Title).
		Suffix("returning id")

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var chatID int64
	if err = pool.QueryRow(ctx, query, args...).Scan(&chatID); err != nil {
		log.Fatalf("failed to insert chat: %v", err)
	}

	log.Printf("inserted chat with id: %d", chatID)

	// добавим информацию о пользователях в чате
	builderUsers := sq.Insert("chat_users").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "user_id")

	for _, userID := range req.Info.UserIds {
		builderUsers = builderUsers.Values(chatID, userID)
	}

	query, args, err = builderUsers.ToSql()
	if err != nil {
		log.Fatalf("failed to build insert_users query: %v", err)
	}

	if _, err = pool.Exec(ctx, query, args...); err != nil {
		log.Fatalf("failed to insert chat users: %v", err)
	}

	log.Printf("inserted chat users: %d", chatID)

	return &desc.CreateChatResponse{
		Id: chatID,
	}, nil
}

// DeleteChatDB удаление чата по идентификатору
func DeleteChatDB(ctx context.Context, req *desc.DeleteChatRequest) (*emptypb.Empty, error) {
	// будем не удалять информацию о чате, а менять статус
	builder := sq.Update("chats").
		PlaceholderFormat(sq.Dollar).
		Set("state", 0).
		Set("deleted_at", time.Now()).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	if _, err = pool.Exec(ctx, query, args...); err != nil {
		log.Fatalf("failed to update chat: %v", err)
	}
	return nil, nil
}

// SendMessageDB отправка сообщения в чат
func SendMessageDB(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	// добавим информацию о чате
	builder := sq.Insert("messages").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "user_id", "body").
		Values(req.Info.ChatId, req.Info.UserId, req.Info.Body).
		Suffix("returning id")

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	if _, err = pool.Exec(ctx, query, args...); err != nil {
		log.Fatalf("failed to insert message: %v", err)
	}
	return nil, nil
}

// GetChatMessagesDB получение списка сообщений чата по идентификатору
// отсортированы по убыванию времени создаения сообщений
// отбирает только активные (state = 1) сообщения
func GetChatMessagesDB(ctx context.Context, req *desc.GetChatMessagesRequest) (*desc.GetChatMessagesResponse, error) {
	// будем не удалять информацию о чате, а менять статус
	builder := sq.Select("id", "user_id", "body", "created_at", "type_id").
		PlaceholderFormat(sq.Dollar).
		From("messages").
		Where(sq.Eq{"chat_id": req.Id, "state": 1}).
		OrderBy("created_at DESC").
		Limit(20)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update chat: %v", err)
	}

	list := make([]*desc.Message, 0)
	var id, userID int64
	var typeID int32
	var body string
	var createdAt time.Time

	for rows.Next() {
		err = rows.Scan(&id, &userID, &body, &createdAt, &typeID)
		if err != nil {
			log.Fatalf("failed to scan message: %v", err)
		}
		list = append(list, &desc.Message{
			Id: id,
			Info: &desc.MessageInfo{
				UserId: userID,
				Body:   body,
			},
			CreatedAt: timestamppb.New(createdAt),
			Type:      typeID,
		})
	}

	return &desc.GetChatMessagesResponse{
		Messages: list,
	}, nil
}
