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

func IsValidStatus(s Status) bool {
	return s > MinStatus && s < MaxStatus
}

type StatusEvent struct {
	UserId    string `json:"userId"`
	Status    Status `json:"newStatus"`
	EventTime int64  `json:"eventTime"`
}
