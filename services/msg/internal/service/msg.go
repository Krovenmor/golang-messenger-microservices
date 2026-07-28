package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	MessageId uuid.UUID
	SenderId  uuid.UUID
	Message   string
	CreatedAt time.Time
}

type Profile struct {
	UserId     uuid.UUID
	Name       string
	PublicKey  string
	PrivateKey string
	KDFSalt    string
}

type MessageService interface {
	NewProfile(ctx context.Context, profile Profile) error
	GetProfile(ctx context.Context, userId uuid.UUID) (*Profile, error)

	CreateNewChat(ctx context.Context, fUser, sUser uuid.UUID) (uuid.UUID, error)
	PostMessage(ctx context.Context, chatId uuid.UUID, msg Message) (uuid.UUID, error)
	GetChatHistory(ctx context.Context, chatId uuid.UUID, fromUserId, fromMsgId uuid.UUID, q int) ([]Message, error)
}
