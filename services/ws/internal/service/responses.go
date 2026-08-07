package service

import (
	"encoding/json"
	"log"
	"net/http"
)

type EventType string

const (
	ResponseType EventType = "response"
)

type ResponseEvent struct {
	Type    EventType `json:"type"`
	Payload any       `json:"payload"`
}

type ResponsePayload struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var (
	StatusOkResp      json.RawMessage
	ErrBadJsonResp    json.RawMessage
	ErrBadStatusResp  json.RawMessage
	ErrBadIdResp      json.RawMessage
	ErrBadRequestResp json.RawMessage
	ErrInternalResp   json.RawMessage
)

func ComputeJsonResponses() {
	var err error
	toResp := func(code int, msg string) json.RawMessage {
		if err != nil {
			return nil
		}
		var rMsg json.RawMessage
		rMsg, err = json.Marshal(ResponseEvent{
			Type: ResponseType,
			Payload: ResponsePayload{
				Code: code,
				Msg:  msg,
			},
		})
		return rMsg
	}

	StatusOkResp = toResp(http.StatusOK, "OK")
	ErrBadJsonResp = toResp(http.StatusBadRequest, "Bad JSON")
	ErrBadStatusResp = toResp(http.StatusBadRequest, "Bad Status")
	ErrBadIdResp = toResp(http.StatusBadRequest, "Bad UUID")
	ErrBadRequestResp = toResp(http.StatusBadRequest, "Bad Request")
	ErrInternalResp = toResp(http.StatusInternalServerError, "Internal Error")

	if err != nil {
		log.Printf("ComputeJsonResponses: Trouble with toResp, err: %q", err.Error())
	}
}
