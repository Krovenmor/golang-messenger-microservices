package queries

import (
	"embed"
)

//go:embed *.sql
var EmbedQueries embed.FS

func GetQuery(queryName string) (string, error) {
	file, err := EmbedQueries.ReadFile(queryName)
	if err != nil {
		return "", err
	}
	return string(file), nil
}
