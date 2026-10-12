package domain

import "strings"

type UserInfo struct {
	Login string
	Email string
}

type UserTokens struct {
	Login   string
	RTokens []string
}

func NormalizeLogin(login string) string {
	return strings.ToLower(login)
}

func NormalizeEmail(email string) string {
	return strings.ToLower(email)
}
