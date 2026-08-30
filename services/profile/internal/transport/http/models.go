package http

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

//
// Requests
//

type NewProfileRequestBody struct {
	Name       string `json:"name" validate:"profile_name"`
	UserName   string `json:"userName" validate:"profile_username"`
	PublicKey  string `json:"pubKey" validate:"profile_pubKey"`
	PrivateKey string `json:"encryptedPrvKey" validate:"profile_prvKey"`
	KDFSalt    string `json:"kdfSalt" validate:"profile_salt"`
	KeyNonce   string `json:"keyNonce" validate:"profile_nonce"`
}

type GetProfilesBatchRequestBody struct {
	Profiles []uuid.UUID `json:"profiles" validate:"profiles_batch_ids"`
}

type ChangeNameRequestBody struct {
	Name string `json:"name" validate:"profile_name"`
}

type ChangeBioRequestBody struct {
	Bio string `json:"bio" validate:"profile_bio"`
}

//
// Responses
//

type PrivateProfileResponseBody struct {
	UserId     uuid.UUID `json:"userId"`
	Name       string    `json:"name"`
	UserName   string    `json:"userName"`
	PublicKey  string    `json:"pubKey"`
	PrivateKey string    `json:"encryptedPrvKey"`
	KDFSalt    string    `json:"kdfSalt"`
	KeyNonce   string    `json:"keyNonce"`
	CreatedAt  time.Time `json:"createdAt"`

	Additional json.RawMessage `json:"additional"`
}

type PublicProfileResponseBody struct {
	UserId    uuid.UUID `json:"userId"`
	Name      string    `json:"name"`
	UserName  string    `json:"userName"`
	PublicKey string    `json:"pubKey"`
	CreatedAt time.Time `json:"createdAt"`

	Additional json.RawMessage `json:"additional"`
}
