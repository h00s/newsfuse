package migrate

import (
	"context"

	"github.com/h00s/newsfuse/app/models"
	"github.com/uptrace/bun"
)

type SeedSources struct{}

func (m SeedSources) Name() string {
	return "seed_sources"
}

func (m SeedSources) Up(db *bun.DB) error {
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

	_, err := db.
		NewInsert().
		Model(&sources).
		Exec(context.Background())

	return err
}

func (m SeedSources) Down(db *bun.DB) error {
	sourceNames := []string{
		"klikni.hr",
		"MojPortal.hr",
		"Radio Daruvar",
		"Index.hr",
		"N1Info.hr",
		"Hacker News",
		"Bug",
		"Telegram",
	}

	_, err := db.
		NewDelete().
		Model((*models.Source)(nil)).
		Where("name IN (?)", bun.In(sourceNames)).
		Exec(context.Background())

	return err
}
