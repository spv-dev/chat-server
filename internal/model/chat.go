package model

import (
	"database/sql"
	"time"
)

// Chat модель чата сервисного слоя
type Chat struct {
	ID        int64
	Info      ChatInfo
	State     int32
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}

// ChatInfo модель информации о чате сервисного слоя
type ChatInfo struct {
	Title   string
	UserIDs []int64
}
