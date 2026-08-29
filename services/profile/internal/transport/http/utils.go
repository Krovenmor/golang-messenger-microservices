package http

import (
	"MyMessenger/pkg/utils"
	"net/http"
)

func recv[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	toRecv, err := utils.Recv[T](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return toRecv, nil
}
