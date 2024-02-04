package db

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/db/migrate"
)

func Migrations() raptor.Migrations {
	return raptor.Migrations{
		1: migrate.AddSource,
		2: migrate.SeedSource,
		3: migrate.AddHeadline,
		4: migrate.AddStory,
	}
}
