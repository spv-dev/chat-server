package repository

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	model "github.com/spv-dev/chat-server/internal/model"
)

// ChatRepository описание методов слоя repo
type ChatRepository interface {
	CreateChat(ctx context.Context, info *model.ChatInfo) (int64, error)
	AddUsersToChat(ctx context.Context, chatID int64, userIDs []int64) error
	DeleteChat(ctx context.Context, id int64) (*emptypb.Empty, error)
	SendMessage(ctx context.Context, info *model.MessageInfo) (*emptypb.Empty, error)
	GetChatMessages(ctx context.Context, id int64) ([]*model.Message, error)
}
