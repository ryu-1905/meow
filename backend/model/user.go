package model

import "time"

// User đại diện cho bản ghi người dùng trong cơ sở dữ liệu.
type User struct {
	ID           int64     `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	HashPassword string    `json:"-" db:"hash_password"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
