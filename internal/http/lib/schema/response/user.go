package response

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64         `json:"id"`
	Username  string        `json:"username"`
	Email     string        `json:"email"`
	Age       sql.NullInt32 `json:"age"`
	Role      string        `json:"role"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type UserShort struct {
	ID       int64         `json:"id"`
	Username string        `json:"username"`
	Email    string        `json:"email"`
	Age      sql.NullInt32 `json:"age"`
	Role     string        `json:"role"`
}
