package models

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Password string    `json:"password" db:"password_hash"`
}
