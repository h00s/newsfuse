// Package db provides database access and management for the newsfuse application.
// This file defines the MigrationsFS function, which is intended to return a filesystem containing the database migration files.
package db

import (
	"io/fs"

	_ "github.com/h00s/newsfuse/db/migrations"
)

func MigrationsFS() fs.FS {
	return nil
}
