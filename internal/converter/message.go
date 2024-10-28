package converter

import (
	"github.com/spv-dev/chat-server/internal/model"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToMessageInfoFromDesc конвертер MessageInfo из API слоя в сервисный слой
func ToMessageInfoFromDesc(info *desc.MessageInfo) model.MessageInfo {
	if info == nil {
		return model.MessageInfo{}
	}
	return model.MessageInfo{
		ChatID: info.ChatId,
		UserID: info.UserId,
		Body:   info.Body,
	}
}

// ToMessagesFromService преобразовывает слайс сообщений из модели в ответ АПИ
func ToMessagesFromService(messages []*model.Message) []*desc.Message {
	var result []*desc.Message
	for _, mess := range messages {
		var updatedAt *timestamppb.Timestamp
		if mess.UpdatedAt != nil {
			updatedAt = timestamppb.New(*mess.UpdatedAt)
		}
		var deletedAt *timestamppb.Timestamp
		if mess.DeletedAt != nil {
			deletedAt = timestamppb.New(*mess.DeletedAt)
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
