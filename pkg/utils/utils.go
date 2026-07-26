package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Send[T any](rw http.ResponseWriter, to_send *T) {
	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	err := json.NewEncoder(rw).Encode(*to_send)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func Recv[T any](req *http.Request, to_rcv *T) error {
	err := json.NewDecoder(req.Body).Decode(&to_rcv)
	if err != nil {
		return fmt.Errorf("Not valid JSON: %q", err.Error())
	}
	return nil
}
