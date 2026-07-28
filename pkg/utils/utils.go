package utils

import (
	"MyMessenger/pkg/jwt"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func Send[T any](w http.ResponseWriter, toSend *T) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(toSend)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Recv[T any](r *http.Request, toRcv *T) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(toRcv)
	if err != nil {
		return fmt.Errorf("Not valid JSON: %w", err)
	}
	return nil
}

func GetUuidFromContext(r *http.Request) (uuid.UUID, error) {
	sUserId := r.Context().Value(jwt.UserIdKey).(string)
	uUserId, err := uuid.Parse(sUserId)
	if err != nil {
		return uUserId, fmt.Errorf("GetUuidFromContext(): trouble with parsing uuid: %w", err)
	}
	return uUserId, nil
}
