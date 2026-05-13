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
	goose.AddMigrationNoTxContext(upSeedSources, downSeedSources)
}

func upSeedSources(ctx context.Context, sqldb *sql.DB) error {
	db := bun.NewDB(sqldb, pgdialect.New())
	sources := models.Sources{
		models.Source{Name: "klikni.hr", TopicID: 1, IsScrapable: true},
		models.Source{Name: "MojPortal.hr", TopicID: 1, IsScrapable: true},
		models.Source{Name: "Radio Daruvar", TopicID: 1, IsScrapable: true},
		models.Source{Name: "Index.hr", TopicID: 2, IsScrapable: true},
		models.Source{Name: "N1Info.hr", TopicID: 2, IsScrapable: true},
		models.Source{Name: "Index.hr", TopicID: 3, IsScrapable: true},
		models.Source{Name: "N1Info.hr", TopicID: 3, IsScrapable: true},
		models.Source{Name: "Hacker News", TopicID: 4, IsScrapable: true},
		models.Source{Name: "Bug", TopicID: 4, IsScrapable: true},
		models.Source{Name: "Telegram", TopicID: 2, IsScrapable: true},
	}
	_, err := db.NewInsert().Model(&sources).Exec(ctx)
	return err
}

func downSeedSources(ctx context.Context, sqldb *sql.DB) error {
	db := bun.NewDB(sqldb, pgdialect.New())
	sourceNames := []string{
		"klikni.hr", "MojPortal.hr", "Radio Daruvar",
		"Index.hr", "N1Info.hr", "Hacker News", "Bug", "Telegram",
	}
	_, err := db.NewDelete().
		Model((*models.Source)(nil)).
		Where("name IN (?)", bun.In(sourceNames)).
		Exec(ctx)
	return err
}
