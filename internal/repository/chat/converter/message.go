package converter

import (
	model "github.com/spv-dev/chat-server/internal/model"
	modelRepo "github.com/spv-dev/chat-server/internal/repository/chat/model"
)

// ToMessagesFromRepo конвертер списка сообщений из БД в сервисный слой
func ToMessagesFromRepo(messages []*modelRepo.Message) []*model.Message {
	var result []*model.Message
	for _, mess := range messages {

		result = append(result, &model.Message{
			ID:    mess.ID,
			State: mess.State,
			Type:  mess.Type,
			Info: model.MessageInfo{
				Body:   mess.Info.Body,
				UserID: mess.Info.UserID,
				ChatID: mess.Info.UserID,
			},
			CreatedAt: mess.CreatedAt,
			UpdatedAt: mess.UpdatedAt,
			DeletedAt: mess.DeletedAt,
		})

	}
	return result
}
