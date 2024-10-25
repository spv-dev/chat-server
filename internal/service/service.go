package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/spv-dev/chat-server/internal/model"
)

// ChatService описание методов сервисного слоя
type ChatService interface {
	CreateChat(ctx context.Context, info *model.ChatInfo) (int64, error)
	DeleteChat(ctx context.Context, id int64) (*emptypb.Empty, error)
	SendMessage(ctx context.Context, info *model.MessageInfo) (*emptypb.Empty, error)
	GetChatMessages(ctx context.Context, id int64) ([]*model.Message, error)
}
