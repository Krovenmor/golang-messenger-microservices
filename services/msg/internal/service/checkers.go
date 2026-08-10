package service

import (
	"fmt"

	"github.com/google/uuid"
)

func (m *MessageServiceImpl) checkProfileName(name string) error {
	lName := len(name)
	if lName < m.conf.MinNameLen {
		return fmt.Errorf("name len must be > %d", m.conf.MinNameLen)
	}
	if lName > m.conf.MaxNameLen {
		return fmt.Errorf("name len must be < %d", m.conf.MaxNameLen)
	}
	return nil
}

func (m *MessageServiceImpl) checkProfileUserName(username string) error {
	lUsName := len(username)
	if lUsName < m.conf.MinUserNameLen {
		return fmt.Errorf("userName len must be > %d", m.conf.MinNameLen)
	}
	if lUsName > m.conf.MaxUserNameLen {
		return fmt.Errorf("userName len must be < %d", m.conf.MaxNameLen)
	}
	return nil
}

func (m *MessageServiceImpl) checkProfileKeys(p *Profile) error {
	lPub, lPrv := len(p.PublicKey), len(p.PrivateKey)
	if lPub < m.conf.MinKeysLen {
		return fmt.Errorf("pubKey is too short, min=%d, got=%d", m.conf.MinKeysLen, lPub)
	}
	if lPrv < m.conf.MinKeysLen {
		return fmt.Errorf("prvKey is too short, min=%d, got=%d", m.conf.MinKeysLen, lPrv)
	}
	if lPub > m.conf.MaxPubKeyLen {
		return fmt.Errorf("pubKey are too big, max=%d, got=%d", m.conf.MaxPubKeyLen, lPub)
	}
	if lPrv > m.conf.MaxPrvKeyLen {
		return fmt.Errorf("prvKey are too big, max=%d, got=%d", m.conf.MaxPrvKeyLen, lPrv)
	}
	return nil
}

func (m *MessageServiceImpl) checkProfileSalt(p *Profile) error {
	lSalt := len(p.KDFSalt)
	if lSalt < m.conf.MinKeysLen {
		return fmt.Errorf("salt are too short")
	}
	if lSalt > m.conf.MaxSaltLen {
		return fmt.Errorf("salt are too big")
	}
	return nil
}

func (m *MessageServiceImpl) checkNonce(nonce string) error {
	lNonce := len(nonce)
	if lNonce < m.conf.MinNonceLen {
		return fmt.Errorf("nonce are too short")
	}
	if lNonce > m.conf.MaxNonceLen {
		return fmt.Errorf("nonce are too big")
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
	err = m.checkNonce(p.KeyNonce)
	if err != nil {
		return err
	}
	if p.UserId == uuid.Nil {
		return fmt.Errorf("nil UserId")
	}
	return nil
}

func (m *MessageServiceImpl) checkMessage(msg *ToPostMessage) error {
	if msg.SenderId == uuid.Nil {
		return fmt.Errorf("Bad Sender UUID")
	}
	if len(msg.ReceiverKey) != m.conf.MsgKeysLen {
		return fmt.Errorf("checkMessage: ReceiverKey != MsgKeysLen, %d != %d", len(msg.ReceiverKey), m.conf.MsgKeysLen)
	}
	if len(msg.SenderKey) != m.conf.MsgKeysLen {
		return fmt.Errorf("checkMessage: SenderKey != MsgKeysLen, %d != %d", len(msg.SenderKey), m.conf.MsgKeysLen)
	}
	if len(msg.Nonce) != m.conf.MsgNonceLen {
		return fmt.Errorf("checkMessage: Nonce != MsgNonceLen, %d != %d", len(msg.Nonce), m.conf.MsgNonceLen)
	}
	return m.checkMessageText(msg.Message)
}

func (m *MessageServiceImpl) checkMessageText(text string) error {
	lMsg := len(text)
	if lMsg < m.conf.MinMsgLen {
		return fmt.Errorf("message len must be > %d", m.conf.MinMsgLen)
	}
	if lMsg > m.conf.MaxMsgLen {
		return fmt.Errorf("message len must be < %d", m.conf.MaxMsgLen)
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
