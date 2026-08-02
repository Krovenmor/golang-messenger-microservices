package service

import (
	"context"

	"github.com/google/uuid"
)

type EventPublisher interface {
	PublishNewChat(ctx context.Context, chatId uuid.UUID, usersTo []uuid.UUID)
	PublishNewMessage(ctx context.Context, chatId, msgId uuid.UUID, usersTo []uuid.UUID)
	PublishMessageWasRedacted(ctx context.Context, chatId, msgId uuid.UUID, usersTo []uuid.UUID)
	PublishMessageWasDeleted(ctx context.Context, chatId, msgId uuid.UUID, usersTo []uuid.UUID)
}
