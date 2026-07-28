package service

import (
	"context"

	"github.com/google/uuid"
)

type MessageRepo interface {
	NewProfile(ctx context.Context, profile Profile) error
	NewChat(ctx context.Context, chatId, fUser, sUser uuid.UUID) error

	PostMessage(ctx context.Context, chatId uuid.UUID, msg Message) error

	GetProfile(ctx context.Context, userId uuid.UUID) (*Profile, error)
	GetChatHistory(ctx context.Context, chatId uuid.UUID, fromId uuid.UUID, q int) ([]Message, error)

	IsProfileInChat(ctx context.Context, userId, chatId uuid.UUID) error
}
