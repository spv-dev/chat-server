package model

import (
	"time"
)

// Chat структура чата
type Chat struct {
	ID        int64      `db:"id"`
	Info      ChatInfo   `db:""`
	State     int32      `db:"state"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at,omitempty"`
	DeletedAt *time.Time `db:"deleted_at,omitempty"`
}

// ChatInfo структура информации о чате
type ChatInfo struct {
	Title   string `db:"title"`
	UserIDs []int64
}
