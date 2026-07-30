package service

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	MessageId  uuid.UUID
	SenderId   uuid.UUID
	Message    *string
	CreatedAt  *time.Time
	IsRedacted bool
	IsDeleted  bool
	RedactedAt *time.Time
}

type Profile struct {
	UserId     uuid.UUID
	Name       string
	UserName   string
	PublicKey  string
	PrivateKey string
	KDFSalt    string
	CreatedAt  time.Time
}

type ChatMember struct {
	UserId   uuid.UUID
	JoinedAt time.Time
}

type ChatInfo struct {
	CreatedAt   time.Time
	ChatMembers []ChatMember
}

type ChatFullInfo struct {
	ChatId    uuid.UUID
	MessageId uuid.UUID
	SenderId  uuid.UUID
	Name      *string
	Message   *string
	CreatedAt *time.Time
}
