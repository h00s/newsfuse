package migrate

import (
	"context"

	"github.com/h00s/newsfuse/app/models"
	"github.com/uptrace/bun"
)

type CreateTopic struct{}

func (m CreateTopic) Name() string {
	return "create_topic_table"
}

func (m CreateTopic) Up(db *bun.DB) error {
	_, err := db.
		NewCreateTable().
		Model((*models.Topic)(nil)).
		Exec(context.Background())
	return err
}

func (m CreateTopic) Down(db *bun.DB) error {
	_, err := db.
		NewDropTable().
		Model((*models.Topic)(nil)).
		IfExists().
		Cascade().
		Exec(context.Background())
	return err
}
