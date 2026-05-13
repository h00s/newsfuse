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
	goose.AddMigrationNoTxContext(upCreateSource, downCreateSource)
}

func upCreateSource(ctx context.Context, sqldb *sql.DB) error {
	db := bun.NewDB(sqldb, pgdialect.New())
	_, err := db.NewCreateTable().Model((*models.Source)(nil)).Exec(ctx)
	return err
}

func downCreateSource(ctx context.Context, sqldb *sql.DB) error {
	db := bun.NewDB(sqldb, pgdialect.New())
	_, err := db.NewDropTable().Model((*models.Source)(nil)).IfExists().Cascade().Exec(ctx)
	return err
}
