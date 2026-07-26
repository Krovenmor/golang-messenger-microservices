package http

type AuthIncomeBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TokensOutcomeBody struct {
	AToken string `json:"accessToken"`
	RToken string `json:"refreshToken"`
}

type TokensUpdateIncome struct {
	RToken string `json:"refreshToken"`
}
