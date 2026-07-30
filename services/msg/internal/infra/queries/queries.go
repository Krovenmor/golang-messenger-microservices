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
	GetProfileId       string
	NewProfile         string
	IsProfileInChat    string
	GetProfileUserName string

	PostMessage   string
	GetMessage    string
	UpdateMessage string
	DeleteMessage string

	NewChat          string
	GetChats         string
	GetChatsExtended string
	GetChatHistory   string
	GetChatInfo      string
	GetChatMembers   string

	GetPrivateChatBetweenTwoPeoples string
}

func GetQueries() (*Queries, error) {
	var err error

	get := func(key string) string {
		if err != nil {
			return ""
		}
		var query string
		query, err = getQuery(key + ".sql")
		return query
	}

	queries := Queries{
		GetChatHistory:                  get("get_chat_history"),
		GetProfileId:                    get("get_profile_id"),
		NewChat:                         get("new_chat"),
		NewProfile:                      get("new_profile"),
		PostMessage:                     get("post_message"),
		IsProfileInChat:                 get("is_profile_in_chat"),
		GetProfileUserName:              get("get_profile_username"),
		GetChats:                        get("get_chats"),
		GetChatInfo:                     get("get_chat_info"),
		GetChatMembers:                  get("get_chat_members"),
		GetPrivateChatBetweenTwoPeoples: get("is_profiles_in_1by1_chat"),
		GetChatsExtended:                get("get_chats_extended"),
		GetMessage:                      get("get_message"),
		UpdateMessage:                   get("update_message"),
		DeleteMessage:                   get("delete_message"),
	}

	if err != nil {
		return nil, err
	}

	return &queries, nil
}
