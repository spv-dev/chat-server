package model

import (
	"database/sql"
	"time"
)

type Chat struct {
	ID        int64        `db:"id"`
	Info      ChatInfo     `db:""`
	State     int32        `db:"state"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type ChatInfo struct {
	Title string `db:"title"`
}
