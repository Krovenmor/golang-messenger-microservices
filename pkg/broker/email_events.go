package broker

const (
	EmailVerificationType   EventType = "email"
	RegistrationAttemptType EventType = "registration_attempt"
)

type EmailVerificationPayload struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type EmailVerificationDTO struct {
	Type    EventType                `json:"type"`
	Payload EmailVerificationPayload `json:"payload"`
}

type RegistrationAttemptPayload struct {
	Email string `json:"email"`
}
