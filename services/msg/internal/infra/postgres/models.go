package postgres

import (
	"MyMessenger/services/msg/internal/service"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	MessageId   uuid.UUID
	SenderId    uuid.UUID
	Message     *string
	SenderKey   string
	ReceiverKey string
	Nonce       string
	CreatedAt   *time.Time
	RedactedAt  *time.Time
	DeletedAt   *time.Time
	ReplyToId   *uuid.UUID
}

type Reaction struct {
	MessageId uuid.UUID
	Emoji     string
	Users     []uuid.UUID
}

func ToServiceMessage(m *Message, r []service.Reaction) service.Message {
	return service.Message{
		MessageId:   m.MessageId,
		SenderId:    m.SenderId,
		Message:     m.Message,
		SenderKey:   m.SenderKey,
		ReceiverKey: m.ReceiverKey,
		Nonce:       m.Nonce,
		CreatedAt:   m.CreatedAt,
		RedactedAt:  m.RedactedAt,
		DeletedAt:   m.DeletedAt,
		ReplyToId:   m.ReplyToId,
		Reactions:   r,
	}
}

func ToServiceReaction(r Reaction) service.Reaction {
	return service.Reaction{
		Emoji: r.Emoji,
		Users: r.Users,
	}
}

func ToServiceReactions(dbReactions []Reaction) []service.Reaction {
	if len(dbReactions) == 0 {
		return []service.Reaction{}
	}
	res := make([]service.Reaction, len(dbReactions))
	for i, r := range dbReactions {
		res[i] = service.Reaction{
			Emoji: r.Emoji,
			Users: r.Users,
		}
	}
	return res
}
