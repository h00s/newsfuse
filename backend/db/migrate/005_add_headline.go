package migrate

import (
	"context"

	"github.com/h00s/newsfuse/app/models"
	"github.com/uptrace/bun"
)

func AddHeadline(db *bun.DB) error {
	_, err := db.
		NewCreateTable().
		Model((*models.Headline)(nil)).
		ForeignKey(`("source_id") REFERENCES "sources" ("id") ON DELETE CASCADE`).
		Exec(context.Background())

	if err != nil {
		return err
	}

	_, err = db.
		NewCreateIndex().
		Model((*models.Headline)(nil)).
		Column("source_id").
		Exec(context.Background())

	if err != nil {
		return err
	}

	_, err = db.
		NewCreateIndex().
		Model((*models.Headline)(nil)).
		Column("published_at").
		Exec(context.Background())

	return err
}
