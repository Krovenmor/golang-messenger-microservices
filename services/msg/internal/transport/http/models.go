package http

import (
	"time"

	"github.com/google/uuid"
)

type ProfileBody struct {
	UserId     uuid.UUID `json:"UserId"`
	Name       string    `json:"Name" validate:"required"`
	UserName   string    `json:"UserName" validate:"required"`
	PublicKey  string    `json:"PubKey" validate:"required"`
	PrivateKey string    `json:"EncryptedPrvKey" validate:"required"`
	KDFSalt    string    `json:"KDFSalt" validate:"required"`
	KeyNonce   string    `json:"KeyNonce" validate:"required"`
	CreatedAt  time.Time `json:"CreatedAt"`
}

type ProfilePublicBody struct {
	UserId    uuid.UUID `json:"UserId"`
	Name      string    `json:"Name"`
	UserName  string    `json:"UserName"`
	PublicKey string    `json:"PubKey"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type NewChatIncomeBody struct {
	UserId uuid.UUID `json:"UserId" validate:"required"`
}

type NewChatResponseBody struct {
	ChatId uuid.UUID `json:"ChatId"`
}

type PostMessageIncomeBody struct {
	Msg string `json:"Message" validate:"required"`
}

type PostMessageResponseBody struct {
	MessageId uuid.UUID `json:"MessageId"`
}
