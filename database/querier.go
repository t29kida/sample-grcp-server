package database

import "context"

//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE
type Querier interface {
	SignUp(context.Context, SignUpParams) (*SignUpResult, error)
	Login(context.Context, LoginParams) (*LoginResult, error)
	CreateSession(context.Context, CreateSessionParams) error
	GetSession(context.Context, GetSessionParams) (*GetSessionResult, error)

	CreateArticle(context.Context, CreateArticleParams) (*CreateArticleResult, error)
	GetArticles(context.Context, GetArticlesParams) (*GetArticlesResult, error)
	GetArticle(context.Context, GetArticleParams) (*GetArticleResult, error)
	UpdateArticle(context.Context, UpdateArticleParams) error
	DeleteArticle(context.Context, DeleteArticleParams) error
}
