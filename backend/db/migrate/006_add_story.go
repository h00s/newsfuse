package migrate

import (
	"context"

	"github.com/h00s/newsfuse/app/models"
	"github.com/uptrace/bun"
)

func AddStory(db *bun.DB) error {
	_, err := db.
		NewCreateTable().
		Model((*models.Story)(nil)).
		Exec(context.Background())
	return err
}
