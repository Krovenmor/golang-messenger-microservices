package service

import (
	"fmt"

	"github.com/google/uuid"
)

func (m *MessageServiceImpl) checkProfileName(name string) error {
	lName := len(name)
	if lName < m.conf.MinNameLen {
		return fmt.Errorf("Name len must be > %d", m.conf.MinNameLen)
	}
	if lName > m.conf.MaxNameLen {
		return fmt.Errorf("Name len must be < %d", m.conf.MaxNameLen)
	}
	return nil
}

func (m *MessageServiceImpl) checkProfileUserName(username string) error {
	lUsName := len(username)
	if lUsName < m.conf.MinUserNameLen {
		return fmt.Errorf("UserName len must be > %d", m.conf.MinNameLen)
	}
	if lUsName > m.conf.MaxUserNameLen {
		return fmt.Errorf("UserName len must be < %d", m.conf.MaxNameLen)
	}
	return nil
}

func (m *MessageServiceImpl) checkProfileKeys(p *Profile) error {
	lPub, lPrv := len(p.PublicKey), len(p.PrivateKey)
	if lPub < m.conf.MinKeysLen {
		return fmt.Errorf("PubKey is too short, min=%d, got=%d", m.conf.MinKeysLen, lPub)
	}
	if lPrv < m.conf.MinKeysLen {
		return fmt.Errorf("PrvKey is too short, min=%d, got=%d", m.conf.MinKeysLen, lPrv)
	}
	if lPub > m.conf.MaxPubKeyLen {
		return fmt.Errorf("PubKey are too big, max=%d, got=%d", m.conf.MaxPubKeyLen, lPub)
	}
	if lPrv > m.conf.MaxPrvKeyLen {
		return fmt.Errorf("PrvKey are too big, max=%d, got=%d", m.conf.MaxPrvKeyLen, lPrv)
	}
	return nil
}

func (m *MessageServiceImpl) checkProfileSalt(p *Profile) error {
	lSalt := len(p.KDFSalt)
	if lSalt < m.conf.MinKeysLen {
		return fmt.Errorf("Salt are too short")
	}
	if lSalt > m.conf.MaxSaltLen {
		return fmt.Errorf("Salt are too big")
	}
	return nil
}

func (m *MessageServiceImpl) checkProfile(p *Profile) error {
	err := m.checkProfileUserName(p.UserName)
	if err != nil {
		return err
	}
	err = m.checkProfileName(p.Name)
	if err != nil {
		return err
	}
	err = m.checkProfileKeys(p)
	if err != nil {
		return err
	}
	err = m.checkProfileSalt(p)
	if err != nil {
		return err
	}
	if p.UserId == uuid.Nil {
		return fmt.Errorf("Nil UserId")
	}
	return nil
}

func (m *MessageServiceImpl) checkMessage(msg *Message) error {
	if msg.Message == nil {
		return fmt.Errorf("No message")
	}
	return m.checkMessageText(*msg.Message)
}

func (m *MessageServiceImpl) checkMessageText(text string) error {
	lMsg := len(text)
	if lMsg < m.conf.MinMsgLen {
		return fmt.Errorf("Message len must be > %d", m.conf.MinMsgLen)
	}
	if lMsg > m.conf.MaxMsgLen {
		return fmt.Errorf("Message len must be < %d", m.conf.MaxMsgLen)
	}
	return nil
}

func (m *MessageServiceImpl) checkQ(q int) error {
	if q > m.conf.MaxQuantity {
		return fmt.Errorf("too big quantity")
	}
	if q < m.conf.MinQuantity {
		return fmt.Errorf("too small quantity")
	}
	return nil
}
