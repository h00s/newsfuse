package db

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/db/migrate"
)

func Migrations() raptor.Migrations {
	return raptor.Migrations{
		1: migrate.AddHeadline,
		2: migrate.AddSource,
		3: migrate.AddStory,
		4: migrate.SeedSource,
	}
}
