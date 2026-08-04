package service

type NewStatusEvent struct {
	NewStatus int `json:"NewStatus"`
}

type Response struct {
	Code int    `json:"Code"`
	Msg  string `json:"Message"`
}
