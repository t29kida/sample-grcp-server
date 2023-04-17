package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

var _ bun.BeforeDropTableHook = (*User)(nil)

func (s *User) BeforeDropTable(ctx context.Context, query *bun.DropTableQuery) error {
	if _, err := query.NewDropTable().IfExists().Table("sessions").Exec(ctx); err != nil {
		return err
	}

	if _, err := query.NewDropTable().IfExists().Table("articles").Exec(ctx); err != nil {
		return err
	}

	return nil
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64        `bun:"id,pk,autoincrement"`
	Email     string       `bun:"email,notnull,unique"`
	Password  string       `bun:"password,notnull"`
	CreatedAt time.Time    `bun:"created_at,notnull,type:timestamp,default:current_timestamp"`
	UpdatedAt time.Time    `bun:"updated_at,notnull,type:timestamp,default:current_timestamp"`
	DeletedAt sql.NullTime `bun:"deleted_at,type:timestamp,soft_delete"`
}

var _ bun.BeforeCreateTableHook = (*Session)(nil)

func (s *Session) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey(`(user_id) REFERENCES users (id) ON DELETE CASCADE`)
	return nil
}

type Session struct {
	bun.BaseModel `bun:"table:sessions,alias:s"`

	AccessToken string    `bun:"access_token,pk"`
	UserID      int64     `bun:"user_id,notnull"`
	CreatedAt   time.Time `bun:"created_at,notnull,type:timestamp,default:current_timestamp"`
	ExpiredAt   time.Time `bun:"expired_at,notnull,type:timestamp"`
}

var _ bun.BeforeCreateTableHook = (*Article)(nil)

func (a *Article) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(user_id) REFERENCES users (id) ON DELETE CASCADE")
	return nil
}

type Article struct {
	bun.BaseModel `bun:"table:articles,alias:a"`

	ID          int64          `bun:"id,pk,autoincrement"`
	UserID      int64          `bun:"user_id,notnull"`
	Title       string         `bun:"title,notnull"`
	Description sql.NullString `bun:"description"`
	Text        string         `bun:"text,notnull,type:text"`
	CreatedAt   time.Time      `bun:"created_at,notnull,type:timestamp,default:current_timestamp"`
	UpdatedAt   time.Time      `bun:"updated_at,notnull,type:timestamp,default:current_timestamp"`
	DeletedAt   sql.NullTime   `bun:"deleted_at,type:timestamp,soft_delete"`
}
