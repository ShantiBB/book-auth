package request

import "database/sql"

type UserCreate struct {
	Username string        `json:"username" validate:"required"`
	Email    string        `json:"email" validate:"required,email"`
	Age      sql.NullInt32 `json:"age"`
	Password string        `json:"password" validate:"required,passwd"`
}

type UserUpdate struct {
	Username string        `json:"username" validate:"required"`
	Age      sql.NullInt32 `json:"age"`
	Email    string        `json:"email" validate:"required,email"`
}
