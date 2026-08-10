package service

import (
	"MyMessenger/pkg/broker"
	"context"

	"github.com/google/uuid"
)

type MessageType int

const (
	TextType MessageType = iota
	BinaryType
)

//go:generate go run github.com/matryer/moq@latest -pkg tests -out ../tests/mocks_test.go . MessageClient Connector Publisher Subscriber

type MessageClient interface {
	GetAllUserChats(ctx context.Context, userId, accessToken string) ([]uuid.UUID, error)
}

type Connector interface {
	Read(ctx context.Context) (MessageType, []byte, error)
	Write(ctx context.Context, msgType MessageType, msg []byte) error
	Ping(ctx context.Context) error
}

type Publisher interface {
	PublishUserStatus(ctx context.Context, status broker.StatusPayload)
	PublishBanEvent(ctx context.Context, payload broker.BanRequestPayload)
}

type Subscriber interface {
	SubscribeOnUserEvents(ctx context.Context, userId string) (<-chan []byte, func(), error)
	SubscribeOnChatEvents(ctx context.Context, chatId string) (<-chan []byte, func(), error)
	SubscribeOnUserStatuses(ctx context.Context, userId string) (<-chan []byte, func(), error)
}
