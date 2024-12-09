package repository

import (
	"context"

	"github.com/spv-dev/chat-server/internal/model"
)

// ChatRepository описание методов repo слоя
type ChatRepository interface {
	CreateChat(ctx context.Context, info *model.ChatInfo) (model.Chat, error)
	DeleteChat(ctx context.Context, id int64) error
	AddUsersToChat(ctx context.Context, chatID int64, userIDs []int64) error
	GetChatMessages(ctx context.Context, id int64, limit uint64, offset uint64) ([]*model.Message, error)
	GetChatInfo(ctx context.Context, id int64) (model.Chat, error)

	SendMessage(ctx context.Context, info *model.MessageInfo) (model.Message, error)
}
