package migrate

import (
	"github.com/go-raptor/raptor/v2"
	"github.com/h00s/newsfuse/app/models"
)

func AddSource(db *raptor.DB) error {
	return db.Migrator().CreateTable(&models.Source{})
}
