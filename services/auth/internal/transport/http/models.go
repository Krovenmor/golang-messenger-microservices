package http

type AuthIncomeBody struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokensOutcomeBody struct {
	AToken string `json:"accessToken"`
	RToken string `json:"refreshToken"`
}

type TokensUpdateIncome struct {
	RToken string `json:"refreshToken" validate:"required"`
}
