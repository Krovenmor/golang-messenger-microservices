package service

import (
	"MyMessenger/services/msg/internal/config"
	"context"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrAlreadyExists = fmt.Errorf("already Exists")
)

type MessageServiceImpl struct {
	repo MessageRepo
	pub  EventPublisher
	conf *config.MessageConfig
}

func NewMessageServiceImpl(repo MessageRepo, pub EventPublisher, conf *config.MessageConfig) *MessageServiceImpl {
	return &MessageServiceImpl{
		repo: repo,
		conf: conf,
		pub:  pub,
	}
}

func (m *MessageServiceImpl) NewProfile(ctx context.Context, profile Profile) error {
	return m.repo.NewProfile(ctx, &profile)
}

func (m *MessageServiceImpl) GetProfileById(ctx context.Context, userId uuid.UUID) (*Profile, error) {
	return m.repo.GetProfileById(ctx, userId)
}

func (m *MessageServiceImpl) GetProfileByUserName(ctx context.Context, username string) (*Profile, error) {
	return m.repo.GetProfileByUserName(ctx, username)
}

func (m *MessageServiceImpl) IsProfileInChat(ctx context.Context, userId, chatId uuid.UUID) error {
	return m.repo.IsProfileInChat(ctx, userId, chatId)
}

func (m *MessageServiceImpl) CreateNewChat(ctx context.Context, fUser, sUser uuid.UUID) (uuid.UUID, error) {
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

func (m *MessageServiceImpl) PostMessage(ctx context.Context, chatId uuid.UUID, msg ToPostMessage) (uuid.UUID, error) {
	msgId, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("trouble with gen UUID for new message: %w", err)
	}
	rMsg := Message{
		MessageId:   msgId,
		SenderId:    msg.UserId,
		Message:     &msg.Message,
		SenderKey:   msg.SenderKey,
		ReceiverKey: msg.ReceiverKey,
		Nonce:       msg.Nonce,
	}
	err = m.repo.NewMessage(ctx, chatId, &rMsg)
	if err != nil {
		return uuid.Nil, err
	}
	m.pub.PublishNewMessage(ctx, chatId, msgId)
	return msgId, nil
}

func (m *MessageServiceImpl) GetMessage(ctx context.Context, chatId, msgId uuid.UUID) (*Message, error) {
	return m.repo.GetMessage(ctx, chatId, msgId)
}

func (m *MessageServiceImpl) RedactMessage(ctx context.Context, chatId, msgId uuid.UUID, msg ToPostMessage) error {
	err := m.repo.RedactMessage(ctx, chatId, msgId, &msg)
	if err != nil {
		return err
	}
	m.pub.PublishMessageWasRedacted(ctx, chatId, msgId)
	return nil
}

func (m *MessageServiceImpl) DelMessage(ctx context.Context, chatId, msgId, userId uuid.UUID) error {
	err := m.repo.DelMessage(ctx, chatId, msgId, userId)
	if err != nil {
		return err
	}
	m.pub.PublishMessageWasDeleted(ctx, chatId, msgId)
	return nil
}

func (m *MessageServiceImpl) GetChats(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	return m.repo.GetChats(ctx, userId)
}

func (m *MessageServiceImpl) GetChatsExtended(ctx context.Context, userId uuid.UUID) ([]ChatFullInfo, error) {
	return m.repo.GetChatsExtended(ctx, userId)
}

func (m *MessageServiceImpl) GetChatInfo(ctx context.Context, chatId uuid.UUID) (*ChatInfo, error) {
	return m.repo.GetChatInfo(ctx, chatId)
}

func (m *MessageServiceImpl) GetChatHistory(ctx context.Context, chatId uuid.UUID, fromUserId, fromMsgId uuid.UUID, q int) ([]Message, error) {
	err := m.checkQ(q)
	if err != nil {
		return []Message{}, err
	}
	return m.repo.GetChatHistory(ctx, chatId, fromMsgId, q)
}
