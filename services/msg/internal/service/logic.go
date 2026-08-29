package service

import (
	"MyMessenger/services/msg/internal/config"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type messageService struct {
	repo MessageRepo
	pub  EventPublisher
	conf *config.MessageConfig
}

func NewMessageService(repo MessageRepo, pub EventPublisher, conf *config.MessageConfig) *messageService {
	return &messageService{
		repo: repo,
		conf: conf,
		pub:  pub,
	}
}

func (m *messageService) NewProfile(ctx context.Context, userId uuid.UUID) error {
	return m.repo.NewProfile(ctx, userId)
}

func (m *messageService) IsProfileInChat(ctx context.Context, userId, chatId uuid.UUID) error {
	return m.repo.IsProfileInChat(ctx, userId, chatId)
}

func (m *messageService) CreateNewChat(ctx context.Context, fUser, sUser uuid.UUID) (uuid.UUID, error) {
	if fUser == sUser {
		return uuid.Nil, fmt.Errorf("can't create chat with fUUID == sUUID")
	}
	chatId, err := m.repo.IsProfilesHaveAPrivateChat(ctx, fUser, sUser)
	if err == nil {
		return chatId, ErrAlreadyExists
	}
	nChatId := uuid.New()
	err = m.repo.NewChat(ctx, nChatId, fUser, sUser)
	if err != nil {
		return uuid.Nil, err
	}
	m.pub.PublishNewChat(ctx, nChatId, []uuid.UUID{sUser})
	return nChatId, nil
}

func (m *messageService) PostMessage(ctx context.Context, chatId uuid.UUID, msg ToPostMessage) (uuid.UUID, error) {
	msgId, err := uuid.NewV7()
	if err != nil {
		log.Printf("trouble with gen UUID for new message: %w", err)
		return uuid.Nil, ErrInternal
	}
	rMsg := Message{
		MessageId:   msgId,
		SenderId:    msg.UserId,
		Message:     &msg.Message,
		SenderKey:   msg.SenderKey,
		ReceiverKey: msg.ReceiverKey,
		Nonce:       msg.Nonce,
		ReplyToId:   msg.ReplyToId,
	}
	err = m.repo.NewMessage(ctx, chatId, &rMsg)
	if err != nil {
		return uuid.Nil, err
	}
	m.pub.PublishNewMessage(ctx, chatId, msgId)
	return msgId, nil
}

func (m *messageService) GetMessage(ctx context.Context, chatId, msgId uuid.UUID) (*Message, error) {
	return m.repo.GetMessage(ctx, chatId, msgId)
}

func (m *messageService) RedactMessage(ctx context.Context, chatId, msgId uuid.UUID, msg ToPostMessage) error {
	err := m.repo.RedactMessage(ctx, chatId, msgId, &msg)
	if err != nil {
		return err
	}
	m.pub.PublishMessageWasRedacted(ctx, chatId, msgId)
	return nil
}

func (m *messageService) DelMessage(ctx context.Context, chatId, msgId, userId uuid.UUID) error {
	err := m.repo.DelMessage(ctx, chatId, msgId, userId)
	if err != nil {
		return err
	}
	m.pub.PublishMessageWasDeleted(ctx, chatId, msgId)
	return nil
}

func (m *messageService) GetChats(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	return m.repo.GetChats(ctx, userId)
}

func (m *messageService) GetChatsExtended(ctx context.Context, userId uuid.UUID) ([]ChatFullInfo, error) {
	return m.repo.GetChatsExtended(ctx, userId)
}

func (m *messageService) GetChatInfo(ctx context.Context, chatId uuid.UUID) (*ChatInfo, error) {
	return m.repo.GetChatInfo(ctx, chatId)
}

func (m *messageService) GetChatHistory(ctx context.Context, chatId uuid.UUID, fromUserId, fromMsgId uuid.UUID, q int) ([]Message, error) {
	err := m.checkQ(q)
	if err != nil {
		return []Message{}, err
	}
	return m.repo.GetChatHistory(ctx, chatId, fromMsgId, q)
}

func (m *messageService) GetSupportedReactions(ctx context.Context) ([]string, error) {
	return m.repo.GetEmojis(ctx)
}

func (m *messageService) PostReaction(ctx context.Context, userId, chatId, msgId uuid.UUID, emoji string) error {
	err := m.repo.NewReaction(ctx, userId, chatId, msgId, emoji)
	if err != nil {
		return err
	}
	m.pub.PublishNewReaction(ctx, chatId, msgId, userId, emoji)
	return nil
}

func (m *messageService) DelReaction(ctx context.Context, userId, chatId, msgId uuid.UUID, emoji string) error {
	err := m.repo.DelReaction(ctx, userId, chatId, msgId, emoji)
	if err != nil {
		return err
	}
	m.pub.PublishDelReaction(ctx, chatId, msgId, userId, emoji)
	return nil
}
