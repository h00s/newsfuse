package migrate

import (
	"context"

	"github.com/h00s/newsfuse/app/models"
	"github.com/uptrace/bun"
)

func AddTopic(db *bun.DB) error {
	_, err := db.
		NewCreateTable().
		Model((*models.Topic)(nil)).
		Exec(context.Background())
	return err
}
