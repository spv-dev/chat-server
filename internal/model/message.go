package model

import (
	"database/sql"
	"time"
)

// Message модель сообщения сервисного слоя
type Message struct {
	ID        int64
	Info      MessageInfo
	State     int32
	Type      int32
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}

// MessageInfo модель информации о сообщении сервисного слоя
type MessageInfo struct {
	ChatID int64
	UserID int64
	Body   string
}
