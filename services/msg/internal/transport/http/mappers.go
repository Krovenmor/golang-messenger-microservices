package http

import (
	"MyMessenger/services/msg/internal/service"

	"github.com/google/uuid"
)

func ToServiceProfile(p *NewProfileRequestBody, userId uuid.UUID) *service.Profile {
	return &service.Profile{
		UserId:     userId,
		Name:       p.Name,
		UserName:   p.UserName,
		PublicKey:  p.PublicKey,
		PrivateKey: p.PrivateKey,
		KDFSalt:    p.KDFSalt,
		KeyNonce:   p.KeyNonce,
	}
}

func ToPrivateProfileBody(p *service.Profile) *PrivateProfileResponseBody {
	return &PrivateProfileResponseBody{
		UserId:     p.UserId,
		Name:       p.Name,
		UserName:   p.UserName,
		PublicKey:  p.PublicKey,
		PrivateKey: p.PrivateKey,
		KDFSalt:    p.KDFSalt,
		CreatedAt:  p.CreatedAt,
		KeyNonce:   p.KeyNonce,
	}
}

func ToPublicProfileBody(p *service.Profile) *PublicProfileResponseBody {
	return &PublicProfileResponseBody{
		UserId:    p.UserId,
		Name:      p.Name,
		UserName:  p.UserName,
		PublicKey: p.PublicKey,
		CreatedAt: p.CreatedAt,
	}
}

func ToServicePostMsg(m *PostMessageRequestBody, userId uuid.UUID) *service.ToPostMessage {
	return &service.ToPostMessage{
		UserId:      userId,
		Message:     m.Msg,
		SenderKey:   m.Skey,
		ReceiverKey: m.Rkey,
		Nonce:       m.Nonce,
	}
}

func FromServiceFullChatInfo(c service.ChatFullInfo) FullChatsInfoResponseBody {
	return FullChatsInfoResponseBody{
		ChatId:      c.ChatId,
		CreatedAt:   c.CreatedAt,
		Members:     c.Members,
		LastMessage: c.LastMessage,
	}
}

func MapSlice[T any, R any](input []T, mapFunc func(T) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = mapFunc(v)
	}
	return result
}

func FromServiceFullChatsInfo(chats []service.ChatFullInfo) []FullChatsInfoResponseBody {
	return MapSlice(chats, FromServiceFullChatInfo)
}

func FromServiceChatInfo(info *service.ChatInfo) *ChatInfoResponseBody {
	return &ChatInfoResponseBody{
		CreatedAt: info.CreatedAt,
		Members:   info.Members,
	}
}

func FromServiceMessage(msg service.Message) MessageResponseBody {
	return MessageResponseBody{
		MessageId:   msg.MessageId,
		SenderId:    msg.SenderId,
		Message:     msg.Message,
		SenderKey:   msg.SenderKey,
		ReceiverKey: msg.ReceiverKey,
		Nonce:       msg.Nonce,
		CreatedAt:   msg.CreatedAt,
		IsRedacted:  msg.IsRedacted,
		IsDeleted:   msg.IsDeleted,
		RedactedAt:  msg.RedactedAt,
	}
}

func FromServiceMessages(msgs []service.Message) []MessageResponseBody {
	return MapSlice(msgs, FromServiceMessage)
}
