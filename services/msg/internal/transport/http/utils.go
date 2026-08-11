package http

import (
	"MyMessenger/pkg/utils"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

func recv[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	toRecv, err := utils.Recv[T](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return toRecv, nil
}

func getChatUuidFromPath(r *http.Request) (uuid.UUID, error) {
	return utils.GetUuidFromPath(r, chatIdKey)
}

func getMsgUuidFromPath(r *http.Request) (uuid.UUID, error) {
	return utils.GetUuidFromPath(r, messageIdKey)
}

func getQueryParam(vals url.Values, key string) (string, error) {
	val := vals.Get(key)
	if val == "" {
		return val, fmt.Errorf("you must provide %q query param", key)
	}
	return val, nil
}

func getUUIDQueryParam(vals url.Values, key string) (uuid.UUID, error) {
	val, err := getQueryParam(vals, key)
	if err != nil {
		return uuid.Nil, err
	}

	valConv, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, fmt.Errorf("not uuid value in %q query param", key)
	}

	return valConv, nil
}

func getIntQueryParam(vals url.Values, key string) (int, error) {
	val, err := getQueryParam(vals, key)
	if err != nil {
		return -1, err
	}

	valConv, err := strconv.Atoi(val)
	if err != nil {
		return -1, fmt.Errorf("not integer value in %q query param", key)
	}

	return valConv, nil
}
