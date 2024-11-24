package entities

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`
	Address   string    `json:"address" db:"address"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
