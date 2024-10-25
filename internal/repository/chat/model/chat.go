package model

import (
	"database/sql"
	"time"
)

// Chat модель БД для чата
type Chat struct {
	ID        int64        `db:"id"`
	Info      ChatInfo     `db:""`
	State     int32        `db:"state"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

// ChatInfo модель информации о чате в БД
type ChatInfo struct {
	Title string `db:"title"`
}
