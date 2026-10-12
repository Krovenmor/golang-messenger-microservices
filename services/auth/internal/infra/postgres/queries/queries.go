package queries

import _ "embed"

var (
	//go:embed add_user.sql
	AddUser string

	//go:embed get_user.sql
	GetUser string

	//go:embed check_taken.sql
	CheckTaken string

	//go:embed get_user_info.sql
	GetUserInfo string

	//go:embed get_user_tokens.sql
	GetUserTokens string

	//go:embed save_refresh.sql
	SaveRefresh string

	//go:embed delete_refresh.sql
	DeleteRefresh string

	//go:embed clear_expired_refresh_tokens.sql
	ClearExpiredRefreshTokens string
)
