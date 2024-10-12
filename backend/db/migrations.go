package db

import (
	"github.com/go-raptor/connector/bun/postgres"
	"github.com/h00s/newsfuse/db/migrate"
)

func Migrations() postgres.Migrations {
	return postgres.Migrations{
		1: migrate.AddTopic,
		2: migrate.SeedTopic,
		3: migrate.AddSource,
		4: migrate.SeedSource,
		5: migrate.AddHeadline,
		6: migrate.AddStory,
	}
}
