package http

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

//
// Requests
//

type NewProfileRequestBody struct {
	Name       string `json:"name" validate:"profile_name"`
	UserName   string `json:"userName" validate:"profile_username"`
	PublicKey  string `json:"pubKey" validate:"profile_pubKey"`
	PrivateKey string `json:"encryptedPrvKey" validate:"profile_prvKey"`
	KDFSalt    string `json:"kdfSalt" validate:"profile_salt"`
	KeyNonce   string `json:"keyNonce" validate:"profile_nonce"`
}

type NewChatRequestBody struct {
	UserId uuid.UUID `json:"userId" validate:"required"`
}

type PostMessageRequestBody struct {
	Msg     string     `json:"message" validate:"message_text"`
	Skey    string     `json:"senderKey" validate:"message_key"`
	Rkey    string     `json:"receiverKey" validate:"message_key"`
	Nonce   string     `json:"nonce" validate:"message_nonce"`
	ReplyId *uuid.UUID `json:"replyToId,omitempty" validate:"omitempty,uuid"`
}

type PostReactionRequestBody struct {
	Emoji string `json:"emoji" validate:"required,min=1,max=64"`
}

//
// Responses
//

type PrivateProfileResponseBody struct {
	UserId     uuid.UUID `json:"userId"`
	Name       string    `json:"name"`
	UserName   string    `json:"userName"`
	PublicKey  string    `json:"pubKey"`
	PrivateKey string    `json:"encryptedPrvKey"`
	KDFSalt    string    `json:"kdfSalt"`
	KeyNonce   string    `json:"keyNonce"`
	CreatedAt  time.Time `json:"createdAt"`
}

type PublicProfileResponseBody struct {
	UserId    uuid.UUID `json:"userId"`
	Name      string    `json:"name"`
	UserName  string    `json:"userName"`
	PublicKey string    `json:"pubKey"`
	CreatedAt time.Time `json:"createdAt"`
}

type NewChatResponseBody struct {
	ChatId uuid.UUID `json:"chatId"`
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

type ReactionResponseBody struct {
	Emoji string      `json:"emoji"`
	Users []uuid.UUID `json:"users"`
}

type MessageResponseBody struct {
	MessageId   uuid.UUID              `json:"messageId"`
	SenderId    uuid.UUID              `json:"senderId"`
	Message     *string                `json:"message"`
	SenderKey   string                 `json:"senderKey"`
	ReceiverKey string                 `json:"receiverKey"`
	Nonce       string                 `json:"nonce"`
	CreatedAt   *time.Time             `json:"createdAt"`
	RedactedAt  *time.Time             `json:"redactedAt"`
	DeletedAt   *time.Time             `json:"deletedAt"`
	ReplyToId   *uuid.UUID             `json:"replyToId"`
	Reactions   []ReactionResponseBody `json:"reactions"`
}
