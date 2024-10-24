package converter

import (
	"github.com/spv-dev/chat-server/internal/model"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUserInfoFromDesc конвертер UserInfo из API слоя в сервисный слой
func ToMessageInfoFromDesc(info *desc.MessageInfo) *model.MessageInfo {
	return &model.MessageInfo{
		ChatID: info.ChatId,
		UserID: info.UserId,
		Body:   info.Body,
	}
}

func ToMessagesFromService(messages []*model.Message) []*desc.Message {
	var result []*desc.Message
	for _, mess := range messages {
		var updatedAt *timestamppb.Timestamp
		if mess.UpdatedAt.Valid {
			updatedAt = timestamppb.New(mess.UpdatedAt.Time)
		}
		var deletedAt *timestamppb.Timestamp
		if mess.DeletedAt.Valid {
			deletedAt = timestamppb.New(mess.DeletedAt.Time)
		}

		result = append(result, &desc.Message{
			Id:    mess.ID,
			State: mess.State,
			Type:  mess.Type,
			Info: &desc.MessageInfo{
				Body:   mess.Info.Body,
				UserId: mess.Info.UserID,
				ChatId: mess.Info.UserID,
			},
			CreatedAt: timestamppb.New(mess.CreatedAt),
			UpdatedAt: updatedAt,
			DeletedAt: deletedAt,
		})

	}
	return result
}
