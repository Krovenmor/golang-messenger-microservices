package http

import (
	"MyMessenger/pkg/utils"
	"errors"
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

func getUuidFromPath(req *http.Request) (uuid.UUID, error) {
	UUID := req.PathValue("uuid")
	if UUID == "" {
		return uuid.Nil, errors.New("Empty UUID")
	}
	cUUID, err := uuid.Parse(UUID)
	if err != nil {
		return uuid.Nil, errors.New("Bad UUID")
	}
	return cUUID, nil
}

func getUUIDQueryParam(vals url.Values, key string) (uuid.UUID, error) {
	val := vals.Get(key)
	if val == "" {
		return uuid.Nil, fmt.Errorf("You must provide %q query param", key)
	}

	valConv, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Not uuid value in %q query param", key)
	}

	return valConv, nil
}

func getIntQueryParam(vals url.Values, key string) (int, error) {
	val := vals.Get(key)
	if val == "" {
		return -1, fmt.Errorf("You must provide %q query param", key)
	}

	valConv, err := strconv.Atoi(val)
	if err != nil {
		return -1, fmt.Errorf("Not integer value in %q query param", key)
	}

	return valConv, nil
}
