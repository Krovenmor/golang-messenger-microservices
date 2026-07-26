package tests

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UserInfo struct {
	login, password string
}

type RepoMock struct {
	commands []string

	repo       map[UserInfo]uuid.UUID
	repoTokens map[uuid.UUID]string
}

func NewMockRepo() *RepoMock {
	return &RepoMock{
		commands:   make([]string, 0),
		repo:       make(map[UserInfo]uuid.UUID),
		repoTokens: make(map[uuid.UUID]string),
	}
}

func (r *RepoMock) AddNewUser(ctx context.Context, userId uuid.UUID, login, password string) error {
	r.commands = append(r.commands, fmt.Sprintf("AddNewUser(): userId=%q, login=%q, password=%q", userId.String(), login, password))
	ui := UserInfo{login: login, password: password}
	r.repo[ui] = userId
	return nil
}
func (r *RepoMock) GetUser(ctx context.Context, login, password string) (uuid.UUID, error) {
	r.commands = append(r.commands, fmt.Sprintf("GetUser(): login=%q, password=%q", login, password))
	ui := UserInfo{login: login, password: password}

	userId, isExists := r.repo[ui]
	if !isExists {
		return uuid.Nil, fmt.Errorf("User not exists: login=%q, password=%q", login, password)
	}

	return userId, nil
}

func (r *RepoMock) SaveRefresh(ctx context.Context, userId uuid.UUID, rToken string) error {
	r.commands = append(r.commands, fmt.Sprintf("SaveRefresh(): userId=%q, rToken=%q", userId.String(), rToken))
	r.repoTokens[userId] = rToken
	return nil
}
func (r *RepoMock) FindRefresh(ctx context.Context, userId uuid.UUID, rToken string) error {
	r.commands = append(r.commands, fmt.Sprintf("FindRefresh(): userId=%q, rToken=%q", userId.String(), rToken))
	token, isExists := r.repoTokens[userId]
	if !isExists {
		return fmt.Errorf("Token not exists: userId=%q", userId.String())
	}
	if rToken != token {
		return fmt.Errorf("Tokens are not equal")
	}
	return nil
}
func (r *RepoMock) DeleteRefresh(ctx context.Context, userId uuid.UUID, rToken string) error {
	r.commands = append(r.commands, fmt.Sprintf("DeleteRefresh(): userId=%q, rToken=%q", userId.String(), rToken))
	token, isExists := r.repoTokens[userId]
	if !isExists {
		return fmt.Errorf("Token not exists: userId=%q", userId.String())
	}
	if rToken != token {
		return fmt.Errorf("Tokens are not equal")
	}
	return nil
}
