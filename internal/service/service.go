package service

import (
	"context"

	"github.com/spv-dev/chat-server/internal/model"
)

// ChatService описание методов сервисного слоя
type ChatService interface {
	CreateChat(ctx context.Context, info *model.ChatInfo) (model.Chat, error)
	DeleteChat(ctx context.Context, id int64) error
	GetChatInfo(ctx context.Context, id int64) (model.Chat, error)

	GetChatMessages(ctx context.Context, id int64, limit uint64, offset uint64) ([]*model.Message, error)

	SendMessage(ctx context.Context, info *model.MessageInfo) (model.Message, error)
}
