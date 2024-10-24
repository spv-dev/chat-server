package model

import (
	"database/sql"
	"time"
)

type Chat struct {
	ID        int64
	Info      ChatInfo
	State     int32
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}

type ChatInfo struct {
	Title   string
	UserIds []int64
}
