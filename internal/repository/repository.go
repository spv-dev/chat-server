package repository

import (
	"context"

	"github.com/spv-dev/chat-server/internal/model"
)

// ChatRepository описание методов repo слоя
type ChatRepository interface {
	CreateChat(ctx context.Context, info *model.ChatInfo) (int64, error)
	DeleteChat(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, info *model.MessageInfo) error
	GetChatMessages(ctx context.Context, id int64) ([]*model.Message, error)
}
