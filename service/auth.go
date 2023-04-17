package service

import (
	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE
type Auther interface {
	CreateAccessToken() (string, error)
}

type Auth struct{}

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) CreateAccessToken() (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", xerrors.Errorf("failed to generate uuid: %v", err)
	}

	return uid.String(), nil
}
