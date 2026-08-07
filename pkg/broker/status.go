package broker

type Status int

const (
	MinStatus Status = iota

	Offline
	Online
	Away
	Typing

	MaxStatus
)

const (
	StatusEvent EventType = "status"
)

func IsValidStatus(s Status) bool {
	return s > MinStatus && s < MaxStatus
}

type StatusPayload struct {
	UserId    string `json:"userId"`
	Status    Status `json:"newStatus"`
	EventTime int64  `json:"eventTime"`
}
