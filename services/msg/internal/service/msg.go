package service

import (
	"context"

	"github.com/google/uuid"
)

type MessageService interface {
	NewProfile(ctx context.Context, profile Profile) error
	GetProfileById(ctx context.Context, userId uuid.UUID) (*Profile, error)
	GetProfileByUserName(ctx context.Context, username string) (*Profile, error)

	CreateNewChat(ctx context.Context, fUser, sUser uuid.UUID) (uuid.UUID, error)
	PostMessage(ctx context.Context, chatId uuid.UUID, msg Message) (uuid.UUID, error)

	GetChatInfo(ctx context.Context, chatId uuid.UUID) (*ChatInfo, error)
	GetChats(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	GetChatHistory(ctx context.Context, chatId uuid.UUID, fromUserId, fromMsgId uuid.UUID, q int) ([]Message, error)
}
