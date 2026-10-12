package domain

import "github.com/google/uuid"

type Registration struct {
	ID           uuid.UUID
	Login        string
	Email        string
	PasswordHash string
	EmailTaken   bool
}
