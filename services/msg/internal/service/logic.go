package service

import (
	"MyMessenger/services/msg/internal/config"
	"context"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrAlreadyExists = fmt.Errorf("Already Exists")
)

type MessageServiceImpl struct {
	repo MessageRepo
	conf *config.MessageConfig
}

func NewMessageServiceImpl(repo MessageRepo, conf *config.MessageConfig) *MessageServiceImpl {
	return &MessageServiceImpl{
		repo: repo,
		conf: conf,
	}
}

func (m *MessageServiceImpl) NewProfile(ctx context.Context, profile Profile) error {
	err := m.checkProfile(&profile)
	if err != nil {
		return err
	}
	return m.repo.NewProfile(ctx, &profile)
}

func (m *MessageServiceImpl) GetProfileById(ctx context.Context, userId uuid.UUID) (*Profile, error) {
	return m.repo.GetProfileById(ctx, userId)
}

func (m *MessageServiceImpl) GetProfileByUserName(ctx context.Context, username string) (*Profile, error) {
	err := m.checkProfileUserName(username)
	if err != nil {
		return nil, err
	}
	return m.repo.GetProfileByUserName(ctx, username)
}

func (m *MessageServiceImpl) CreateNewChat(ctx context.Context, fUser, sUser uuid.UUID) (uuid.UUID, error) {
	if fUser == sUser {
		return uuid.Nil, fmt.Errorf("Can't create chat with fUUID == sUUID")
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
	return nChatId, nil
}

func (m *MessageServiceImpl) PostMessage(ctx context.Context, chatId uuid.UUID, msg Message) (uuid.UUID, error) {
	msgId, err := uuid.NewV7()
	if err != nil {
		return msgId, fmt.Errorf("Trouble with gen UUID for new message: %w", err)
	}
	msg.MessageId = msgId
	return msgId, m.repo.PostMessage(ctx, chatId, &msg)
}

func (m *MessageServiceImpl) GetChats(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	return m.repo.GetChats(ctx, userId)
}

func (m *MessageServiceImpl) GetChatInfo(ctx context.Context, chatId uuid.UUID) (*ChatInfo, error) {
	return m.repo.GetChatInfo(ctx, chatId)
}

func (m *MessageServiceImpl) GetChatHistory(ctx context.Context, chatId uuid.UUID, fromUserId, fromMsgId uuid.UUID, q int) ([]Message, error) {
	err := m.checkFromAndQ(fromMsgId, q)
	if err != nil {
		return []Message{}, err
	}
	err = m.repo.IsProfileInChat(ctx, fromUserId, chatId)
	if err != nil {
		return []Message{}, fmt.Errorf("Can't access other chats")
	}
	return m.repo.GetChatHistory(ctx, chatId, fromMsgId, q)
}
