package queries

import (
	"embed"
)

//go:embed *.sql
var EmbedQueries embed.FS

func getQuery(queryName string) (string, error) {
	file, err := EmbedQueries.ReadFile(queryName)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

type Queries struct {
	GetChatHistory  string
	GetProfile      string
	NewChat         string
	NewProfile      string
	PostMessage     string
	IsProfileInChat string
}

func GetQueries() (*Queries, error) {
	GetChatHistory, err := getQuery("get_chat_history.sql")
	if err != nil {
		return nil, err
	}
	GetProfile, err := getQuery("get_profile.sql")
	if err != nil {
		return nil, err
	}
	NewChat, err := getQuery("new_chat.sql")
	if err != nil {
		return nil, err
	}
	NewProfile, err := getQuery("new_profile.sql")
	if err != nil {
		return nil, err
	}
	PostMessage, err := getQuery("post_message.sql")
	if err != nil {
		return nil, err
	}
	IsProfileInChat, err := getQuery("is_profile_in_chat.sql")
	if err != nil {
		return nil, err
	}
	return &Queries{
		GetChatHistory:  GetChatHistory,
		GetProfile:      GetProfile,
		NewChat:         NewChat,
		NewProfile:      NewProfile,
		PostMessage:     PostMessage,
		IsProfileInChat: IsProfileInChat,
	}, nil
}
