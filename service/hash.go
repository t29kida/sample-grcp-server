package service

import (
	"github.com/alexedwards/argon2id"
	"golang.org/x/xerrors"
)

//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE
type Hasher interface {
	CreateHash(string) (string, error)
	CompareHash(string, string) (bool, error)
}

type Hash struct {
	params *argon2id.Params
}

func NewHash() *Hash {
	return &Hash{
		params: argon2id.DefaultParams,
	}
}

func (hash *Hash) CreateHash(password string) (string, error) {
	h, err := argon2id.CreateHash(password, hash.params)
	if err != nil {
		return "", xerrors.Errorf("failed to create hash: %v", err)
	}

	return h, nil
}

func (hash *Hash) CompareHash(password, hashed string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hashed)
	if err != nil {
		return false, xerrors.Errorf("failed to compare hash: %v", err)
	}

	return match, nil
}
