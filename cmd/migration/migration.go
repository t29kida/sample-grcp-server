package migration

import (
	"context"
	"database/sql"
	"log"
	"time"

	cfg "sample-grpc-server/config"
	"sample-grpc-server/database"
	"sample-grpc-server/database/model"

	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"golang.org/x/xerrors"
)

var Cmd = &cobra.Command{
	Use:          "migration",
	SilenceUsage: false,
	Run: func(cmd *cobra.Command, args []string) {
		Start()
	},
}

func openDB(ctx context.Context) (*sql.DB, error) {
	cfg.LoadConfig()

	loc := time.FixedZone("Local", 9*60*60)

	dns := database.FormatDNS(loc)

	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, xerrors.Errorf("failed to open database connection: %v", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, xerrors.Errorf("failed to verify database connection: %v", err)
	}

	return db, nil
}

func Start() {
	log.Println("started migration")

	ctx := context.Background()

	sqlDB, err := openDB(ctx)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}

	db := bun.NewDB(sqlDB, mysqldialect.New())

	if err := db.ResetModel(ctx,
		(*model.User)(nil),
		(*model.Session)(nil),
		(*model.Article)(nil),
	); err != nil {
		log.Fatalf("failed to reset tables: %v", err)
	}

	log.Println("completed migration")
}
