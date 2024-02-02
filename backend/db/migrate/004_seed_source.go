package migrate

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/models"
)

func SeedSource(db *raptor.DB) error {
	sources := models.Sources{
		models.Source{Name: "klikni.hr"},
		models.Source{Name: "MojPortal.hr"},
		models.Source{Name: "Radio Daruvar"},
	}
	return db.CreateInBatches(&sources, 3).Error
}
