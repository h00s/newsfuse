package services

import (
	"log"

	"github.com/h00s/newsfuse/api/config"
	"github.com/h00s/newsfuse/api/db"
	"github.com/h00s/newsfuse/api/scrapers"
)

type Services struct {
	DB       *db.Database
	Logger   *log.Logger
	Scrapers []scrapers.Scraper
}

func NewServices(config *config.Config, logger *log.Logger) *Services {
	db := db.NewDatabase(&config.Database)
	if err := db.Connect(); err != nil {
		logger.Fatal(err)
	}
	if err := db.Migrate(); err != nil {
		logger.Fatal(err)
	}

	return &Services{
		DB:       db,
		Logger:   logger,
		Scrapers: scrapers.NewScrapers(),
	}
}

func (s *Services) Close() {
	s.DB.Close()
}
