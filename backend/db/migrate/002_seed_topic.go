package migrate

import (
	"context"

	"github.com/h00s/newsfuse/app/models"
	"github.com/uptrace/bun"
)

func SeedTopic(db *bun.DB) error {
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
