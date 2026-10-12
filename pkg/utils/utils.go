package utils

import (
	"MyMessenger/pkg/jwt"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Send[T any](w http.ResponseWriter, toSend *T) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(toSend)
	if err != nil {
		zap.L().Error("failed to encode response", zap.Error(err))
	}
}

func SendWithStatus[T any](w http.ResponseWriter, toSend *T, status int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(toSend)
	if err != nil {
		zap.L().Error("failed to encode response", zap.Error(err))
	}
}

var Validator = validator.New()

func Recv[T any](r *http.Request) (*T, error) {
	var val T
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&val)
	if err != nil {
		return nil, fmt.Errorf("not valid JSON: %w", err)
	}
	err = Validator.Struct(&val)
	if err != nil {
		return nil, fmt.Errorf("not valid JSON: %w", err)
	}
	return &val, nil
}

func GetUuidFromContext(r *http.Request) (uuid.UUID, error) {
	sUserId := r.Context().Value(jwt.UserIdKey).(string)
	uUserId, err := uuid.Parse(sUserId)
	if err != nil {
		return uUserId, fmt.Errorf("parse user id from context: %w", err)
	}
	return uUserId, nil
}

func GetUuidFromPath(r *http.Request, key string) (uuid.UUID, error) {
	UUID := r.PathValue(key)
	if UUID == "" {
		return uuid.Nil, fmt.Errorf("empty %q UUID", key)
	}
	cUUID, err := uuid.Parse(UUID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("bad %q UUID", key)
	}
	return cUUID, nil
}

func Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func MapSlice[T any, R any](input []T, mapFunc func(T) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = mapFunc(v)
	}
	return result
}
