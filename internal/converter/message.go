package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/spv-dev/chat-server/internal/model"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

// ToMessageInfoFromDesc конвертер MessageInfo из API слоя в сервисный слой
func ToMessageInfoFromDesc(info *desc.MessageInfo) *model.MessageInfo {
	if info == nil {
		return &model.MessageInfo{}
	}
	return &model.MessageInfo{
		ChatID: info.ChatId,
		UserID: info.UserId,
		Body:   info.Body,
	}
}

// ToMessageInfoFromDesc конвертер MessageInfo из API слоя в сервисный слой
func ToMessageFromModel(info *model.Message) *desc.Message {
	if info == nil {
		return &desc.Message{}
	}
	var updatedAt *timestamppb.Timestamp
	if info.UpdatedAt != nil {
		updatedAt = timestamppb.New(*info.UpdatedAt)
	}
	var deletedAt *timestamppb.Timestamp
	if info.DeletedAt != nil {
		deletedAt = timestamppb.New(*info.DeletedAt)
	}

	return &desc.Message{
		Id: info.ID,
		Info: &desc.MessageInfo{
			ChatId: info.Info.ChatID,
			UserId: info.Info.UserID,
			Body:   info.Info.Body,
		},
		State:     info.State,
		Type:      info.Type,
		CreatedAt: timestamppb.New(info.CreatedAt),
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
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
				ChatId: mess.Info.ChatID,
			},
			CreatedAt: timestamppb.New(mess.CreatedAt),
			UpdatedAt: updatedAt,
			DeletedAt: deletedAt,
		})

	}
	return result
}
