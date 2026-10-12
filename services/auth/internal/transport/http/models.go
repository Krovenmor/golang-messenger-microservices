package http

//
//	Requests
//

type StartRegistrationRequestBody struct {
	Login    string `json:"login" validate:"auth_login"`
	Password string `json:"password" validate:"auth_password"`
	Email    string `json:"email" validate:"required,email"`
}

type ConfirmRegistrationRequestBody struct {
	RegistrationId string `json:"registrationId" validate:"required,uuid"`
	Code           string `json:"code" validate:"required,len=6,numeric"`
}

type ResendCodeRequestBody struct {
	RegistrationId string `json:"registrationId" validate:"required,uuid"`
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

type StartRegistrationResponseBody struct {
	RegistrationId string `json:"registrationId"`
	ExpiresIn      int    `json:"expiresIn"`
	RetryAfter     int    `json:"retryAfter"`
}

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
