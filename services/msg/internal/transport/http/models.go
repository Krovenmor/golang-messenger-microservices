package http

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProfileBody struct {
	UserId     uuid.UUID `json:"userId"`
	Name       string    `json:"name" validate:"required"`
	UserName   string    `json:"userName" validate:"required"`
	PublicKey  string    `json:"pubKey" validate:"required"`
	PrivateKey string    `json:"encryptedPrvKey" validate:"required"`
	KDFSalt    string    `json:"kdfSalt" validate:"required"`
	KeyNonce   string    `json:"keyNonce" validate:"required"`
	CreatedAt  time.Time `json:"createdAt"`
}

type ProfilePublicBody struct {
	UserId    uuid.UUID `json:"userId"`
	Name      string    `json:"name"`
	UserName  string    `json:"userName"`
	PublicKey string    `json:"pubKey"`
	CreatedAt time.Time `json:"createdAt"`
}

type NewChatRequestBody struct {
	UserId uuid.UUID `json:"userId" validate:"required"`
}

type NewChatResponseBody struct {
	ChatId uuid.UUID `json:"chatId"`
}

type PostMessageRequestBody struct {
	Msg   string `json:"message" validate:"required"`
	Skey  string `json:"senderKey" validate:"required"`
	Rkey  string `json:"receiverKey" validate:"required"`
	Nonce string `json:"nonce" validate:"required"`
}

type ChangeMessageRequestBody struct {
	Msg string `json:"message" validate:"required"`
}

type PostMessageResponseBody struct {
	MessageId uuid.UUID `json:"messageId"`
}

type FullChatsInfoResponseBody struct {
	ChatId      uuid.UUID       `json:"chatId"`
	CreatedAt   time.Time       `json:"createdAt"`
	Members     json.RawMessage `json:"members"`
	LastMessage json.RawMessage `json:"lastMessage"`
}

type ChatInfoResponseBody struct {
	CreatedAt time.Time       `json:"createdAt"`
	Members   json.RawMessage `json:"members"`
}

type MessageResponseBody struct {
	MessageId   uuid.UUID  `json:"messageId"`
	SenderId    uuid.UUID  `json:"senderId"`
	Message     *string    `json:"message"`
	SenderKey   string     `json:"senderKey"`
	ReceiverKey string     `json:"receiverKey"`
	Nonce       string     `json:"nonce"`
	CreatedAt   *time.Time `json:"createdAt"`
	IsRedacted  bool       `json:"isRedacted"`
	IsDeleted   bool       `json:"isDeleted"`
	RedactedAt  *time.Time `json:"redactedAt"`
}
