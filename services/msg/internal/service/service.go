package service

import (
	"context"

	"github.com/google/uuid"
)

type MessageService interface {
	NewProfile(ctx context.Context, profile Profile) error
	GetProfileById(ctx context.Context, userId uuid.UUID) (*Profile, error)
	GetProfileByUserName(ctx context.Context, username string) (*Profile, error)
	IsProfileInChat(ctx context.Context, userId, chatId uuid.UUID) error

	AddAvatarToProfile(ctx context.Context, userId uuid.UUID, avatarId uuid.UUID) error
	DelAvatarFromProfile(ctx context.Context, userId uuid.UUID, avatarId uuid.UUID) error

	PostMessage(ctx context.Context, chatId uuid.UUID, msg ToPostMessage) (uuid.UUID, error)
	GetMessage(ctx context.Context, chatId, msgId uuid.UUID) (*Message, error)
	RedactMessage(ctx context.Context, chatId, msgId uuid.UUID, msg ToPostMessage) error
	DelMessage(ctx context.Context, chatId, msgId, userId uuid.UUID) error

	CreateNewChat(ctx context.Context, fUser, sUser uuid.UUID) (uuid.UUID, error)
	GetChatInfo(ctx context.Context, chatId uuid.UUID) (*ChatInfo, error)
	GetChats(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	GetChatsExtended(ctx context.Context, userId uuid.UUID) ([]ChatFullInfo, error)
	GetChatHistory(ctx context.Context, chatId uuid.UUID, fromUserId, fromMsgId uuid.UUID, q int) ([]Message, error)

	GetSupportedReactions(ctx context.Context) ([]string, error)
	PostReaction(ctx context.Context, userId, chatId, msgId uuid.UUID, emoji string) error
	DelReaction(ctx context.Context, userId, chatId, msgId uuid.UUID, emoji string) error
}
