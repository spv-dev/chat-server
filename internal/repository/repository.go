package repository

import (
	"context"

	"github.com/spv-dev/chat-server/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, info *model.ChatInfo) (int64, error)
	DeleteChat(ctx context.Context, id int64) (*emptypb.Empty, error)
	SendMessage(ctx context.Context, info *model.MessageInfo) (*emptypb.Empty, error)
	GetChatMessages(ctx context.Context, id int64) ([]*model.Message, error)
}
