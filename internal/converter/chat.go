package converter

import (
	"github.com/spv-dev/chat-server/internal/model"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

// ToChatInfoFromDesc конвертер ChatInfo из API слоя в сервисный слой
func ToChatInfoFromDesc(info *desc.ChatInfo) model.ChatInfo {
	if info == nil {
		return model.ChatInfo{}
	}
	return model.ChatInfo{
		Title:   info.Title,
		UserIDs: info.UserIds,
	}
}
