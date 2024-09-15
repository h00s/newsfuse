package migrate

import (
	"context"

	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/app/models"
)

func AddHeadline(db *raptor.DB) error {
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
