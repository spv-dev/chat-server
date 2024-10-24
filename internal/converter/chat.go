package converter

import (
	"github.com/spv-dev/chat-server/internal/model"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

// ToUserInfoFromDesc конвертер UserInfo из API слоя в сервисный слой
func ToChatInfoFromDesc(info *desc.ChatInfo) *model.ChatInfo {
	return &model.ChatInfo{
		Title:   info.Title,
		UserIds: info.UserIds,
	}
}
