package http

import (
	"MyMessenger/pkg/utils"
	"MyMessenger/services/msg/internal/service"

	"github.com/google/uuid"
)

func ToServicePostMsg(m *PostMessageRequestBody, userId uuid.UUID) *service.ToPostMessage {
	return &service.ToPostMessage{
		UserId:      userId,
		Message:     m.Msg,
		SenderKey:   m.Skey,
		ReceiverKey: m.Rkey,
		Nonce:       m.Nonce,
		ReplyToId:   m.ReplyId,
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

func FromServiceFullChatsInfo(chats []service.ChatFullInfo) []FullChatsInfoResponseBody {
	return utils.MapSlice(chats, FromServiceFullChatInfo)
}

func FromServiceChatInfo(info *service.ChatInfo) *ChatInfoResponseBody {
	return &ChatInfoResponseBody{
		CreatedAt: info.CreatedAt,
		Members:   info.Members,
	}
}

func FromServiceReaction(reaction service.Reaction) ReactionResponseBody {
	return ReactionResponseBody{
		Emoji: reaction.Emoji,
		Users: reaction.Users,
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
		RedactedAt:  msg.RedactedAt,
		DeletedAt:   msg.DeletedAt,
		ReplyToId:   msg.ReplyToId,
		Reactions:   utils.MapSlice(msg.Reactions, FromServiceReaction),
	}
}

func FromServiceMessages(msgs []service.Message) []MessageResponseBody {
	return utils.MapSlice(msgs, FromServiceMessage)
}
