package service

import (
	"context"

	"github.com/google/uuid"
)

type MessageRepo interface {
	NewProfile(ctx context.Context, userId uuid.UUID) error

	IsProfileInChat(ctx context.Context, userId, chatId uuid.UUID) error
	IsProfilesHaveAPrivateChat(ctx context.Context, userIdF, userIdS uuid.UUID) (uuid.UUID, error)

	NewChat(ctx context.Context, chatId, fUser, sUser uuid.UUID) error
	GetChats(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	GetChatsExtended(ctx context.Context, userId uuid.UUID) ([]ChatFullInfo, error)
	GetChatInfo(ctx context.Context, chatId uuid.UUID) (*ChatInfo, error)
	GetChatMembers(ctx context.Context, chatId uuid.UUID) ([]ChatMember, error)
	GetChatHistory(ctx context.Context, chatId uuid.UUID, fromId uuid.UUID, q int) ([]Message, error)

	NewMessage(ctx context.Context, chatId uuid.UUID, msg *Message) error
	GetMessage(ctx context.Context, chatId, msgId uuid.UUID) (*Message, error)
	RedactMessage(ctx context.Context, chatId, msgId uuid.UUID, msg *ToPostMessage) error
	DelMessage(ctx context.Context, chatId, msgId, userId uuid.UUID) error

	GetEmojis(ctx context.Context) ([]string, error)
	NewReaction(ctx context.Context, userId, chatId, msgId uuid.UUID, emoji string) error
	DelReaction(ctx context.Context, userId, chatId, msgId uuid.UUID, emoji string) error
}

type EventPublisher interface {
	PublishNewChat(ctx context.Context, chatId uuid.UUID, usersTo []uuid.UUID)
	PublishNewMessage(ctx context.Context, chatId, msgId uuid.UUID)
	PublishMessageWasRedacted(ctx context.Context, chatId, msgId uuid.UUID)
	PublishMessageWasDeleted(ctx context.Context, chatId, msgId uuid.UUID)

	PublishNewReaction(ctx context.Context, chatId, msgId, userId uuid.UUID, emoji string)
	PublishDelReaction(ctx context.Context, chatId, msgId, userId uuid.UUID, emoji string)
}
