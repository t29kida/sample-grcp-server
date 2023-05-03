package database

import (
	"database/sql"
	"os"
	"time"

	cfg "sample-grpc-server/config"

	"github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
	"golang.org/x/xerrors"
)

func NewDatabase() (*bun.DB, error) {
	loc := time.FixedZone("Local", 9*60*60)

	dns := FormatDNS(loc)

	sqlDB, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, xerrors.Errorf("failed to open connection: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, xerrors.Errorf("failed to verify connection: %v", err)
	}

	db := bun.NewDB(sqlDB, mysqldialect.New())
	if os.Getenv("ENV") == "development" {
		db.AddQueryHook(
			bundebug.NewQueryHook(
				bundebug.WithEnabled(true),
			),
		)
	}

	return db, nil
}

func FormatDNS(loc *time.Location) string {
	c := mysql.Config{
		User:      cfg.Cfg.GetDBUser(),
		Passwd:    cfg.Cfg.GetDBPassword(),
		Net:       "tcp",
		Addr:      cfg.Cfg.GetDBAddr(),
		DBName:    cfg.Cfg.GetDBName(),
		Loc:       loc,
		ParseTime: true,
	}

	return c.FormatDSN()
}
