package model

import (
	"time"
)

// Message структура сообщения
type Message struct {
	ID        int64       `db:"id"`
	Info      MessageInfo `db:""`
	State     int32       `db:"state"`
	Type      int32       `db:"type_id"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt *time.Time  `db:"updated_at,omitempty"`
	DeletedAt *time.Time  `db:"deleted_at,omitempty"`
}

// MessageInfo структура информации о сообщении
type MessageInfo struct {
	ChatID int64  `db:"chat_id"`
	UserID int64  `db:"user_id"`
	Body   string `db:"body"`
}
