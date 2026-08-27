package queries

import (
	"MyMessenger/pkg/repo"
	"embed"
)

//go:embed *.sql
var EmbedQueries embed.FS

type Queries struct {
	NewData           string
	GetAvailableSpace string
	DelData           string

	NewProfile     string
	GetProfileInfo string
	GetProfileData string
}

func GetQueries() (*Queries, error) {
	r := repo.NewQueriesReader(EmbedQueries)

	queries := Queries{
		NewData:           r.GetQuery("new_data"),
		GetAvailableSpace: r.GetQuery("get_available_space"),
		DelData:           r.GetQuery("del_data"),

		NewProfile:     r.GetQuery("new_profile"),
		GetProfileInfo: r.GetQuery("get_profile_info"),
		GetProfileData: r.GetQuery("get_profile_data"),
	}

	return &queries, r.Err()
}
