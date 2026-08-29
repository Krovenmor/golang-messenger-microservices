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
	NewProfile string

	IsProfileInChat string

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

	GetEmojis    string
	GetReactions string
	NewReaction  string
	DelReaction  string
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
		NewProfile: get("new_profile"),

		GetChatHistory:                  get("get_chat_history"),
		NewChat:                         get("new_chat"),
		PostMessage:                     get("post_message"),
		IsProfileInChat:                 get("is_profile_in_chat"),
		GetChats:                        get("get_chats"),
		GetChatInfo:                     get("get_chat_info"),
		GetChatMembers:                  get("get_chat_members"),
		GetPrivateChatBetweenTwoPeoples: get("is_profiles_in_1by1_chat"),
		GetChatsExtended:                get("get_chats_extended"),
		GetMessage:                      get("get_message"),
		UpdateMessage:                   get("update_message"),
		DeleteMessage:                   get("delete_message"),

		GetEmojis:    get("get_emojis"),
		GetReactions: get("get_reactions"),
		NewReaction:  get("new_reaction"),
		DelReaction:  get("delete_reaction"),
	}

	if err != nil {
		return nil, err
	}

	return &queries, nil
}
