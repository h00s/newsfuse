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
	goose.AddMigrationNoTxContext(upCreateHeadline, downCreateHeadline)
}

func upCreateHeadline(ctx context.Context, sqldb *sql.DB) error {
	db := bun.NewDB(sqldb, pgdialect.New())
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewCreateTable().
			Model((*models.Headline)(nil)).
			ForeignKey(`("source_id") REFERENCES "sources" ("id") ON DELETE CASCADE`).
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			Model((*models.Headline)(nil)).
			Column("source_id").
			Exec(ctx); err != nil {
			return err
		}
		_, err := tx.NewCreateIndex().
			Model((*models.Headline)(nil)).
			Column("published_at").
			Exec(ctx)
		return err
	})
}

func downCreateHeadline(ctx context.Context, sqldb *sql.DB) error {
	db := bun.NewDB(sqldb, pgdialect.New())
	_, err := db.NewDropTable().
		Model((*models.Headline)(nil)).
		IfExists().Cascade().
		Exec(ctx)
	return err
}
