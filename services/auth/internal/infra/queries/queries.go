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
	AddUser       string
	GetUser       string
	SaveRefresh   string
	FindRefresh   string
	DelRefresh    string
	ClrExpRefresh string
}

func GetQueries() (*Queries, error) {
	addUser, err := getQuery("add_user.sql")
	if err != nil {
		return nil, err
	}
	getUser, err := getQuery("get_user.sql")
	if err != nil {
		return nil, err
	}
	saveRefresh, err := getQuery("save_refresh.sql")
	if err != nil {
		return nil, err
	}
	findRefresh, err := getQuery("find_refresh.sql")
	if err != nil {
		return nil, err
	}
	delRefresh, err := getQuery("delete_refresh.sql")
	if err != nil {
		return nil, err
	}
	clrRefresh, err := getQuery("clear_expired_refresh_tokens.sql")
	if err != nil {
		return nil, err
	}
	return &Queries{
		AddUser:       addUser,
		GetUser:       getUser,
		SaveRefresh:   saveRefresh,
		FindRefresh:   findRefresh,
		DelRefresh:    delRefresh,
		ClrExpRefresh: clrRefresh,
	}, nil
}
