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
	AddUser         string
	GetUser         string
	SaveRefresh     string
	FindRefresh     string
	DelRefresh      string
	ClrExpRefresh   string
	CheckUserExists string
	GetUserInfo     string
}

func GetQueries() (*Queries, error) {
	var err error
	get := func(key string) string {
		if err != nil {
			return ""
		}
		var val string
		val, err = getQuery(key + ".sql")
		return val
	}

	q := Queries{
		AddUser:         get("add_user"),
		GetUser:         get("get_user"),
		SaveRefresh:     get("save_refresh"),
		FindRefresh:     get("find_refresh"),
		DelRefresh:      get("delete_refresh"),
		ClrExpRefresh:   get("clear_expired_refresh_tokens"),
		CheckUserExists: get("check_user_exists"),
		GetUserInfo:     get("get_user_info"),
	}

	return &q, err
}
