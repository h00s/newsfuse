package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/h00s/newsfuse/app/models"
)

func init() {
	goose.AddMigrationNoTxContext(upSeedTopics, downSeedTopics)
}

func upSeedTopics(ctx context.Context, sqldb *sql.DB) error {
	db := bun.NewDB(sqldb, pgdialect.New())
	topics := models.Topics{
		models.Topic{Name: "BBŽ"},
		models.Topic{Name: "Hrvatska"},
		models.Topic{Name: "Svijet"},
		models.Topic{Name: "Tech"},
	}
	_, err := db.NewInsert().Model(&topics).Exec(ctx)
	return err
}

func downSeedTopics(ctx context.Context, sqldb *sql.DB) error {
	db := bun.NewDB(sqldb, pgdialect.New())
	_, err := db.NewDelete().
		Model((*models.Topic)(nil)).
		Where("name IN (?)", bun.In([]string{"BBŽ", "Hrvatska", "Svijet", "Tech"})).
		Exec(ctx)
	return err
}
