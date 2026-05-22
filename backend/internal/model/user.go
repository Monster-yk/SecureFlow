package model

import (
	"time"

	"github.com/google/uuid"
)

const(
	RoleAdmin = "admin"
	RoleAnalyst = "analyst"
)

type User struct{
	ID uuid.UUID  `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	Role string `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Updatedat time.Time `json:"updated_at" db:"updated_at"`
}