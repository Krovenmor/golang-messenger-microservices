package broker

const (
	NewProfileEvent EventType = "new:profile"
)

type ProfileEventDTO struct {
	Type    EventType      `json:"type"`
	Payload ProfilePayload `json:"payload"`
}

type ProfilePayload struct {
	UserId string `json:"userId"`
}
