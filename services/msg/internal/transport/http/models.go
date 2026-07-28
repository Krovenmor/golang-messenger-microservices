package http

import (
	"github.com/google/uuid"
)

type ProfileBody struct {
	Name       string `json:"UserName"`
	PublicKey  string `json:"PubKey"`
	PrivateKey string `json:"EncryptedPrvKey"`
	KDFSalt    string `json:"KDFSalt"`
}

type NewChatIncomeBody struct {
	UserId uuid.UUID `json:"UserId"`
}

type NewChatResponseBody struct {
	ChatId uuid.UUID `json:"ChatId"`
}

type PostMessageIncomeBody struct {
	Msg string `json:"Message"`
}

type PostMessageResponseBody struct {
	MessageId uuid.UUID `json:"MessageId"`
}
