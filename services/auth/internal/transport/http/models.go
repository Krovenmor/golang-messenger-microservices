package http

//
//	Requests
//

type AuthRequestBody struct {
	Login    string `json:"login" validate:"auth_login"`
	Password string `json:"password" validate:"auth_password"`
}

type TokensUpdateRequestBody struct {
	RToken string `json:"refreshToken" validate:"refresh_token"`
}

//
//	Responses
//

type TokensResponseBody struct {
	AToken string `json:"accessToken"`
	RToken string `json:"refreshToken"`
}
