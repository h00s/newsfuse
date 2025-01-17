package migrate

import (
	"context"

	"github.com/h00s/newsfuse/app/models"
	"github.com/uptrace/bun"
)

type SeedTopics struct{}

func (m SeedTopics) Name() string {
	return "seed_topics"
}

func (m SeedTopics) Up(db *bun.DB) error {
	topics := models.Topics{
		models.Topic{Name: "BBŽ"},
		models.Topic{Name: "Hrvatska"},
		models.Topic{Name: "Svijet"},
		models.Topic{Name: "Tech"},
	}

	_, err := db.
		NewInsert().
		Model(&topics).
		Exec(context.Background())

	return err
}

func (m SeedTopics) Down(db *bun.DB) error {
	_, err := db.
		NewDelete().
		Model((*models.Topic)(nil)).
		Where("name IN (?)", bun.In([]string{"BBŽ", "Hrvatska", "Svijet", "Tech"})).
		Exec(context.Background())

	return err
}
