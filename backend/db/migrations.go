package db

import (
	"io/fs"

	_ "github.com/h00s/newsfuse/db/migrations"
)

func MigrationsFS() fs.FS {
	return nil
}
