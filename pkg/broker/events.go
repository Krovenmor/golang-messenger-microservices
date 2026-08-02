package broker

type EventType string

const (
	NewChatType         EventType = "newChat"
	NewMessageType      EventType = "newMessage"
	MessageRedactedType EventType = "messageRedacted"
	MessageDeletedType  EventType = "messageDeleted"
)

type Event struct {
	Type    EventType `json:"type"`
	Payload any       `json:"payload"`
}

type NewChatPayload struct {
	ChatId string `json:"chatId"`
}

type NewMessagePayload struct {
	ChatId string `json:"chatId"`
	MsgId  string `json:"msgId"`
}

type MessageRedactedPayload struct {
	ChatId string `json:"chatId"`
	MsgId  string `json:"msgId"`
}

type MessageDeletedPayload struct {
	ChatId string `json:"chatId"`
	MsgId  string `json:"msgId"`
}
