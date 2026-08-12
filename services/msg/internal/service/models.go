package service

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ToPostMessage struct {
	UserId      uuid.UUID
	Message     string
	SenderKey   string
	ReceiverKey string
	Nonce       string
	ReplyToId   *uuid.UUID
}

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

type Profile struct {
	UserId     uuid.UUID
	Name       string
	UserName   string
	PublicKey  string
	PrivateKey string
	KDFSalt    string
	KeyNonce   string
	CreatedAt  time.Time
}

type ChatMember struct {
	UserId   uuid.UUID
	Name     string
	JoinedAt time.Time
}

type ChatInfo struct {
	CreatedAt time.Time
	Members   json.RawMessage
}

type ChatFullInfo struct {
	ChatId      uuid.UUID
	CreatedAt   time.Time
	Members     json.RawMessage
	LastMessage json.RawMessage
}
