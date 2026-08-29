package service

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	UserId     uuid.UUID
	Name       string
	UserName   string
	PublicKey  string
	PrivateKey string
	KDFSalt    string
	KeyNonce   string
	CreatedAt  time.Time
	Additional json.RawMessage
}
