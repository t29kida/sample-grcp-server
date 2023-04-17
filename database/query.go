package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"sample-grpc-server/database/model"

	"github.com/uptrace/bun"
	"golang.org/x/xerrors"
)

type Query struct {
	db *bun.DB
}

func NewQuery(db *bun.DB) *Query {
	return &Query{
		db: db,
	}
}

type SignUpParams struct {
	Email    string
	Password string
}

type SignUpResult struct {
	UserID int64
}

func (q *Query) SignUp(ctx context.Context, p SignUpParams) (*SignUpResult, error) {
	user := model.User{
		Email:    p.Email,
		Password: p.Password,
	}

	result, err := q.db.NewInsert().Model(&user).Exec(ctx)
	if err != nil {
		return nil, xerrors.Errorf("新規ユーザー登録: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, xerrors.Errorf("failed to get last inserted id: %v", err)
	}

	return &SignUpResult{
		UserID: userID,
	}, nil
}

type LoginParams struct {
	Email string
}

type LoginResult struct {
	UserID   int64
	Password string
}

func (q *Query) Login(ctx context.Context, p LoginParams) (*LoginResult, error) {
	user := new(model.User)

	err := q.db.NewSelect().
		ColumnExpr("id").
		ColumnExpr("password").
		TableExpr("users").
		Where("email = ?", p.Email).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(ctx, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, xerrors.Errorf("ユーザー取得: %w", err)
		} else {
			return nil, xerrors.Errorf("ユーザー取得: %v", err)
		}
	}

	return &LoginResult{UserID: user.ID, Password: user.Password}, nil
}

type CreateSessionParams struct {
	AccessToken string
	UserID      int64
}

func (q *Query) CreateSession(ctx context.Context, p CreateSessionParams) error {
	session := model.Session{
		AccessToken: p.AccessToken,
		UserID:      p.UserID,
		ExpiredAt:   time.Now().Add(time.Hour * 24),
	}

	if _, err := q.db.NewInsert().Model(&session).Exec(ctx); err != nil {
		return xerrors.Errorf("failed to create session: %v", err)
	}

	return nil
}

type GetSessionParams struct {
	AccessToken string
}

type GetSessionResult struct {
	ID int64
}

func (q *Query) GetSession(ctx context.Context, p GetSessionParams) (*GetSessionResult, error) {
	session := new(model.Session)

	err := q.db.NewSelect().Column("user_id").
		Table("sessions").
		Where("access_token = ?", p.AccessToken).
		Where("expired_at > CURRENT_TIMESTAMP").
		Scan(ctx, session)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, xerrors.Errorf("failed to get session: %w", err)
		}
		return nil, xerrors.Errorf("failed to get session: %v", err)
	}

	return &GetSessionResult{ID: session.UserID}, nil
}

type CreateArticleParams struct {
	UserID      int64
	Title       string
	Description sql.NullString
	Text        string
}

type CreateArticleResult struct {
	ArticleID int64
}

func (q *Query) CreateArticle(ctx context.Context, p CreateArticleParams) (*CreateArticleResult, error) {
	article := model.Article{
		UserID:      p.UserID,
		Title:       p.Title,
		Description: p.Description,
		Text:        p.Text,
	}

	result, err := q.db.NewInsert().Model(&article).Exec(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to create article: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, xerrors.Errorf("failed to get last inserted id: %v", err)
	}

	return &CreateArticleResult{ArticleID: id}, nil
}

type GetArticlesParams struct {
	UserID int64
}

type GetArticlesResult struct {
	Articles []model.Article
}

func (q *Query) GetArticles(ctx context.Context, p GetArticlesParams) (*GetArticlesResult, error) {
	var articles []model.Article

	err := q.db.NewSelect().
		Column("id").
		Column("title").
		Column("description").
		Column("text").
		Table("articles").
		Where("user_id = ?", p.UserID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Scan(ctx, &articles)

	if err != nil {
		return nil, xerrors.Errorf("failed to get articles: %v", err)
	}

	return &GetArticlesResult{Articles: articles}, nil
}

type GetArticleParams struct {
	ArticleID int64
	UserID    int64
}

type GetArticleResult struct {
	Article model.Article
}

func (q *Query) GetArticle(ctx context.Context, p GetArticleParams) (*GetArticleResult, error) {
	var article model.Article

	err := q.db.NewSelect().
		Column("id").
		Column("title").
		Column("description").
		Column("text").
		Table("articles").
		Where("id = ?", p.ArticleID).
		Where("user_id = ?", p.UserID).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(ctx, &article)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, xerrors.Errorf("article not found: %w", err)
		}
		return nil, xerrors.Errorf("failed to get articles: %v", err)
	}

	return &GetArticleResult{Article: article}, nil
}

type UpdateArticleParams struct {
	ArticleID   int64
	UserID      int64
	Title       string
	Description sql.NullString
	Text        string
}

func (q *Query) UpdateArticle(ctx context.Context, p UpdateArticleParams) error {
	_, err := q.db.NewUpdate().
		Table("articles").
		Set("title = ?", p.Title).
		Set("description = ?", p.Description).
		Set("text = ?", p.Text).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", p.ArticleID).
		Where("user_id = ?", p.UserID).
		Returning("NULL").
		Exec(ctx)

	if err != nil {
		return xerrors.Errorf("failed to update article: %v", err)
	}

	return nil
}

type DeleteArticleParams struct {
	ArticleID int64
	UserID    int64
}

func (q *Query) DeleteArticle(ctx context.Context, p DeleteArticleParams) error {
	_, err := q.db.NewUpdate().
		Table("articles").
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", p.ArticleID).
		Where("user_id = ?", p.UserID).
		Returning("NULL").
		Exec(ctx)

	if err != nil {
		return xerrors.Errorf("failed to delete article: %v", err)
	}

	return nil
}
