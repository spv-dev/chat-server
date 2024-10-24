package model

import (
	"database/sql"
	"time"
)

type Message struct {
	ID        int64        `db:"id"`
	Info      MessageInfo  `db:""`
	State     int32        `db:"state"`
	Type      int32        `db:"type"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type MessageInfo struct {
	ChatID int64  `db:"chat_id"`
	UserID int64  `db:"user_id"`
	Body   string `db:"body"`
}
