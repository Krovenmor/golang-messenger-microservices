package queries

import (
	"MyMessenger/pkg/repo"
	"embed"
)

//go:embed *.sql
var EmbedQueries embed.FS

type Queries struct {
	NewProfile string
	DelProfile string

	GetProfileUserName string
	GetProfileId       string
	GetProfilesId      string

	AddAvatarPhoto string
	DelAvatarPhoto string
}

func GetQueries() (*Queries, error) {
	r := repo.NewQueriesReader(EmbedQueries)

	queries := Queries{
		NewProfile: r.GetQuery("new_profile"),
		DelProfile: r.GetQuery("del_profile"),

		GetProfileId:       r.GetQuery("get_profile_id"),
		GetProfilesId:      r.GetQuery("get_profiles_id"),
		GetProfileUserName: r.GetQuery("get_profile_username"),

		AddAvatarPhoto: r.GetQuery("add_avatar_photo"),
		DelAvatarPhoto: r.GetQuery("del_avatar_photo"),
	}

	return &queries, r.Err()
}
