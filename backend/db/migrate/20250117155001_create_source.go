package migrate

import (
	"context"

	"github.com/h00s/newsfuse/app/models"
	"github.com/uptrace/bun"
)

type CreateSource struct{}

func (m CreateSource) Name() string {
	return "create_source_table"
}

func (m CreateSource) Up(db *bun.DB) error {
	_, err := db.
		NewCreateTable().
		Model((*models.Source)(nil)).
		Exec(context.Background())
	return err
}

func (m CreateSource) Down(db *bun.DB) error {
	_, err := db.
		NewDropTable().
		Model((*models.Source)(nil)).
		IfExists().
		Cascade().
		Exec(context.Background())
	return err
}
