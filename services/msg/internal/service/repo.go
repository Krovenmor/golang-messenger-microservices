package service

import (
	"context"

	"github.com/google/uuid"
)

type MessageRepo interface {
	NewProfile(ctx context.Context, profile *Profile) error
	GetProfileById(ctx context.Context, userId uuid.UUID) (*Profile, error)
	GetProfileByUserName(ctx context.Context, username string) (*Profile, error)
	IsProfileInChat(ctx context.Context, userId, chatId uuid.UUID) error
	IsProfilesHaveAPrivateChat(ctx context.Context, userIdF, userIdS uuid.UUID) (uuid.UUID, error)

	NewChat(ctx context.Context, chatId, fUser, sUser uuid.UUID) error
	GetChats(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	GetChatsExtended(ctx context.Context, userId uuid.UUID) ([]ChatFullInfo, error)
	GetChatInfo(ctx context.Context, chatId uuid.UUID) (*ChatInfo, error)
	GetChatMembers(ctx context.Context, chatId uuid.UUID) ([]ChatMember, error)
	GetChatMembersIdsExcept(ctx context.Context, chatId, exceptUserId uuid.UUID) ([]uuid.UUID, error)
	GetChatHistory(ctx context.Context, chatId uuid.UUID, fromId uuid.UUID, q int) ([]Message, error)

	NewMessage(ctx context.Context, chatId uuid.UUID, msg *Message) error
	GetMessage(ctx context.Context, chatId, msgId uuid.UUID) (*Message, error)
	RedactMessage(ctx context.Context, chatId, msgId, userId uuid.UUID, newText string) error
	DelMessage(ctx context.Context, chatId, msgId, userId uuid.UUID) error
}
