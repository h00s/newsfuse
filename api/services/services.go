package services

import (
	"log"

	"github.com/h00s/newsfuse/api/config"
	"github.com/h00s/newsfuse/api/db"
)

type Services struct {
	DB     *db.Database
	Logger *log.Logger
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
		DB:     db,
		Logger: logger,
	}
}

func (s *Services) Close() {
	s.DB.Close()
}
