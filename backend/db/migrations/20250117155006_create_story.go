package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/h00s/newsfuse/app/models"
)

func init() {
	goose.AddMigrationNoTxContext(upCreateStory, downCreateStory)
}

func upCreateStory(ctx context.Context, sqldb *sql.DB) error {
	db := bun.NewDB(sqldb, pgdialect.New())
	_, err := db.NewCreateTable().Model((*models.Story)(nil)).Exec(ctx)
	return err
}

func downCreateStory(ctx context.Context, sqldb *sql.DB) error {
	db := bun.NewDB(sqldb, pgdialect.New())
	_, err := db.NewDropTable().Model((*models.Story)(nil)).IfExists().Cascade().Exec(ctx)
	return err
}
