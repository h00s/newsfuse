package migrate

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/models"
)

func SeedTopic(db *raptor.DB) error {
	topics := models.Topics{
		models.Topic{Name: "BBŽ"},
		models.Topic{Name: "Hrvatska"},
	}
	return db.CreateInBatches(&topics, 2).Error
}
