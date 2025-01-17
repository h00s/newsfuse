package migrate

import (
	"context"

	"github.com/h00s/newsfuse/app/models"
	"github.com/uptrace/bun"
)

type CreateHeadline struct{}

func (m CreateHeadline) Name() string {
	return "create_headline_table"
}

func (m CreateHeadline) Up(db *bun.DB) error {
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

func (m CreateHeadline) Down(db *bun.DB) error {
	_, err := db.
		NewDropTable().
		Model((*models.Headline)(nil)).
		IfExists().
		Cascade().
		Exec(context.Background())
	return err
}
