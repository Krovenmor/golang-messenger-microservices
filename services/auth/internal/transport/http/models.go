package http

//
//	Requests
//

type SendCodeRequestBody struct {
	Email string `json:"email" validate:"required,email"`
}

type RegisterRequestBody struct {
	Login    string `json:"login" validate:"auth_login"`
	Password string `json:"password" validate:"auth_password"`
	Email    string `json:"email" validate:"required,email"`
	Code     string `json:"code" validate:"required,len=6"`
}

type LoginRequestBody struct {
	Login    string `json:"login" validate:"auth_login"`
	Password string `json:"password" validate:"auth_password"`
}

type TokensUpdateRequestBody struct {
	RToken string `json:"refreshToken" validate:"refresh_token"`
}

//
//	Responses
//

type SendCodeResponseBody struct {
	ExpiresIn  int `json:"expiresIn"`
	RetryAfter int `json:"retryAfter"`
}

type TokensResponseBody struct {
	AToken string `json:"accessToken"`
	RToken string `json:"refreshToken"`
}

type AccountInfoResponseBody struct {
	Login string `json:"login"`
	Email string `json:"email"`
}
