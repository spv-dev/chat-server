package converter

import (
	"github.com/spv-dev/chat-server/internal/model"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToChatInfoFromDesc конвертер ChatInfo из API слоя в сервисный слой
func ToChatInfoFromDesc(info *desc.ChatInfo) *model.ChatInfo {
	if info == nil {
		return &model.ChatInfo{}
	}
	return &model.ChatInfo{
		Title:   info.Title,
		UserIDs: info.UserIds,
	}
}

// ToChatFromModel конвертер Chat в API слой из сервисного слоя
func ToChatFromModel(info *model.Chat) *desc.Chat {
	if info == nil {
		return &desc.Chat{}
	}
	var updatedAt *timestamppb.Timestamp
	if info.UpdatedAt != nil {
		updatedAt = timestamppb.New(*info.UpdatedAt)
	}
	var deletedAt *timestamppb.Timestamp
	if info.DeletedAt != nil {
		deletedAt = timestamppb.New(*info.DeletedAt)
	}

	return &desc.Chat{
		Id:    info.ID,
		State: info.State,
		Info: &desc.ChatInfo{
			Title: info.Info.Title,
		},
		CreatedAt: timestamppb.New(info.CreatedAt),
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}
}
