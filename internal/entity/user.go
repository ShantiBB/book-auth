package entity

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64         `json:"id"`
	Username     string        `json:"username"`
	Email        string        `json:"email"`
	Age          sql.NullInt32 `json:"age"`
	PasswordHash string        `json:"password_hash"`
	Role         string        `json:"role"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}
