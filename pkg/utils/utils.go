package utils

import (
	"MyMessenger/pkg/jwt"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func Send[T any](w http.ResponseWriter, toSend *T) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(toSend)
	if err != nil {
		log.Printf("Send: trouble with encoding, err: %q", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

func SendWithStatus[T any](w http.ResponseWriter, toSend *T, status int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(toSend)
	if err != nil {
		log.Printf("SendWithStatus: trouble with encoding, err: %q", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

var validate = validator.New()

func Recv[T any](r *http.Request) (*T, error) {
	var val T
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&val)
	if err != nil {
		return nil, fmt.Errorf("Not valid JSON: %w", err)
	}
	err = validate.Struct(&val)
	if err != nil {
		return nil, fmt.Errorf("Not valid JSON: %w", err)
	}
	return &val, nil
}

func GetUuidFromContext(r *http.Request) (uuid.UUID, error) {
	sUserId := r.Context().Value(jwt.UserIdKey).(string)
	uUserId, err := uuid.Parse(sUserId)
	if err != nil {
		err := fmt.Errorf("GetUuidFromContext(): trouble with parsing uuid: %w", err)
		log.Print(err.Error())
		return uUserId, err
	}
	return uUserId, nil
}

func GetUuidFromPath(r *http.Request, key string) (uuid.UUID, error) {
	UUID := r.PathValue(key)
	if UUID == "" {
		return uuid.Nil, fmt.Errorf("Empty %q UUID", key)
	}
	cUUID, err := uuid.Parse(UUID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Bad %q UUID", key)
	}
	return cUUID, nil
}
