package migrate

import (
	"context"

	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/app/models"
)

func AddStory(db *raptor.DB) error {
	_, err := db.
		NewCreateTable().
		Model((*models.Story)(nil)).
		Exec(context.Background())
	return err
}
