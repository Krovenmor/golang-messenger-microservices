package service

import (
	"MyMessenger/services/msg/internal/config"
	"context"
	"fmt"

	"github.com/google/uuid"
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

func (m *MessageServiceImpl) checkProfile(p *Profile) error {
	lName := len(p.Name)
	if lName < m.conf.MinNameLen {
		return fmt.Errorf("Name len must be > %d", m.conf.MinNameLen)
	}
	if lName > m.conf.MaxNameLen {
		return fmt.Errorf("Name len must be < %d", m.conf.MaxNameLen)
	}
	if p.UserId == uuid.Nil {
		return fmt.Errorf("Nil UserId")
	}
	lPub, lPrv, lSalt := len(p.PublicKey), len(p.PrivateKey), len(p.KDFSalt)
	if lPub < m.conf.MinKeysLen || lPrv < m.conf.MinKeysLen {
		return fmt.Errorf("Keys are too short")
	}
	if lSalt < m.conf.MinKeysLen {
		return fmt.Errorf("Salt are too short")
	}
	if lPub > m.conf.MaxPubKeyLen || lPrv > m.conf.MaxPrvKeyLen {
		return fmt.Errorf("Keys are too big")
	}
	if lSalt > m.conf.MaxSaltLen {
		return fmt.Errorf("Salt are too big")
	}
	return nil
}

func (m *MessageServiceImpl) NewProfile(ctx context.Context, profile Profile) error {
	err := m.checkProfile(&profile)
	if err != nil {
		return err
	}
	return m.repo.NewProfile(ctx, profile)
}

func (m *MessageServiceImpl) GetProfile(ctx context.Context, userId uuid.UUID) (*Profile, error) {
	return m.repo.GetProfile(ctx, userId)
}

func (m *MessageServiceImpl) CreateNewChat(ctx context.Context, fUser, sUser uuid.UUID) (uuid.UUID, error) {
	if fUser == sUser {
		return uuid.Nil, fmt.Errorf("Can't create chat with fUUID == sUUID")
	}
	nChatId := uuid.New()
	err := m.repo.NewChat(ctx, nChatId, fUser, sUser)
	if err != nil {
		return uuid.Nil, err
	}
	return nChatId, nil
}

func (m *MessageServiceImpl) checkMsg(msg *Message) error {
	lMsg := len(msg.Message)
	if lMsg < m.conf.MinMsgLen {
		return fmt.Errorf("Message len must be > %d", m.conf.MinMsgLen)
	}
	if lMsg > m.conf.MaxMsgLen {
		return fmt.Errorf("Message len must be < %d", m.conf.MaxMsgLen)
	}
	return nil
}

func (m *MessageServiceImpl) PostMessage(ctx context.Context, chatId uuid.UUID, msg Message) (uuid.UUID, error) {
	msgId, err := uuid.NewV7()
	if err != nil {
		return msgId, fmt.Errorf("Trouble with gen UUID for new message: %w", err)
	}
	msg.MessageId = msgId
	return msgId, m.repo.PostMessage(ctx, chatId, msg)
}

func (m *MessageServiceImpl) checkFromAndQ(fromId uuid.UUID, q int) error {
	if fromId == uuid.Nil {
		return fmt.Errorf("from can't be Nil")
	}
	if q > m.conf.MaxQuantity {
		return fmt.Errorf("too big quantity")
	}
	if q < m.conf.MinQuantity {
		return fmt.Errorf("too small quantity")
	}
	return nil
}

func (m *MessageServiceImpl) GetChatHistory(ctx context.Context, chatId uuid.UUID, fromUserId, fromMsgId uuid.UUID, q int) ([]Message, error) {
	err := m.checkFromAndQ(fromMsgId, q)
	if err != nil {
		return []Message{}, err
	}
	err = m.repo.IsProfileInChat(ctx, fromUserId, fromMsgId)
	if err != nil {
		return []Message{}, fmt.Errorf("Can't access other chats")
	}
	return m.repo.GetChatHistory(ctx, chatId, fromMsgId, q)
}
