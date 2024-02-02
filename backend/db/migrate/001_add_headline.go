package migrate

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/models"
)

func AddHeadline(db *raptor.DB) error {
	return db.Migrator().CreateTable(&models.Headline{})
}
