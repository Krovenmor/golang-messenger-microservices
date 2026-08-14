package broker

const (
	NewChatType         EventType = "newChat"
	NewMessageType      EventType = "newMessage"
	MessageRedactedType EventType = "messageRedacted"
	MessageDeletedType  EventType = "messageDeleted"

	NewReactionType EventType = "newReaction"
	DelReactionType EventType = "delReaction"
)

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

type NewReactionPayload struct {
	ChatId string `json:"chatId"`
	MsgId  string `json:"msgId"`
	UserId string `json:"userId"`
	Emoji  string `json:"emoji"`
}

type DelReactionPayload struct {
	ChatId string `json:"chatId"`
	MsgId  string `json:"msgId"`
	UserId string `json:"userId"`
	Emoji  string `json:"emoji"`
}
