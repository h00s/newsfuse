package db

import (
	"github.com/go-raptor/connectors/bun/postgres"
	"github.com/h00s/newsfuse/db/migrate"
)

func Migrations() postgres.Migrations {
	return postgres.Migrations{
		"20250117155001": &migrate.CreateSource{},
		"20250117155002": &migrate.SeedSources{},
		"20250117155003": &migrate.CreateTopic{},
		"20250117155004": &migrate.SeedTopics{},
		"20250117155005": &migrate.CreateHeadline{},
		"20250117155006": &migrate.CreateStory{},
	}
}
