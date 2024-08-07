package migrate

import (
	"github.com/go-raptor/raptor/v2"
	"github.com/h00s/newsfuse/app/models"
)

func SeedSource(db *raptor.DB) error {
	sources := models.Sources{
		models.Source{Name: "klikni.hr", TopicID: 1, IsScrapable: true},
		models.Source{Name: "MojPortal.hr", TopicID: 1, IsScrapable: true},
		models.Source{Name: "Radio Daruvar", TopicID: 1, IsScrapable: true},
		models.Source{Name: "Index.hr", TopicID: 2, IsScrapable: true},
		models.Source{Name: "N1Info.hr", TopicID: 2, IsScrapable: true},
		models.Source{Name: "Index.hr", TopicID: 3, IsScrapable: true},
		models.Source{Name: "N1Info.hr", TopicID: 3, IsScrapable: true},
		models.Source{Name: "Hacker News", TopicID: 4, IsScrapable: true},
		models.Source{Name: "Bug", TopicID: 4, IsScrapable: true},
		models.Source{Name: "Telegram", TopicID: 2, IsScrapable: true},
	}
	return db.CreateInBatches(&sources, 10).Error
}
