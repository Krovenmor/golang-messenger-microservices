package service

import (
	"context"

	"github.com/google/uuid"
)

type MessageRepo interface {
	NewProfile(ctx context.Context, profile *Profile) error
	NewChat(ctx context.Context, chatId, fUser, sUser uuid.UUID) error

	PostMessage(ctx context.Context, chatId uuid.UUID, msg *Message) error

	GetProfileById(ctx context.Context, userId uuid.UUID) (*Profile, error)
	GetProfileByUserName(ctx context.Context, username string) (*Profile, error)
	GetChats(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	GetChatInfo(ctx context.Context, chatId uuid.UUID) (*ChatInfo, error)
	GetChatHistory(ctx context.Context, chatId uuid.UUID, fromId uuid.UUID, q int) ([]Message, error)

	IsProfileInChat(ctx context.Context, userId, chatId uuid.UUID) error
	IsProfilesHaveAPrivateChat(ctx context.Context, userIdF, userIdS uuid.UUID) (uuid.UUID, error)
}
