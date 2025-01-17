package migrate

import (
	"context"

	"github.com/h00s/newsfuse/app/models"
	"github.com/uptrace/bun"
)

type CreateStory struct{}

func (m CreateStory) Name() string {
	return "create_story_table"
}

func (m CreateStory) Up(db *bun.DB) error {
	_, err := db.
		NewCreateTable().
		Model((*models.Story)(nil)).
		Exec(context.Background())
	return err
}

func (m CreateStory) Down(db *bun.DB) error {
	_, err := db.
		NewDropTable().
		Model((*models.Story)(nil)).
		IfExists().
		Cascade().
		Exec(context.Background())
	return err
}
