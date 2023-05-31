package services

import (
	"log"

	"github.com/h00s/newsfuse/api/config"
	"github.com/h00s/newsfuse/api/db"
	"github.com/h00s/newsfuse/api/robot"
)

type Services struct {
	DB     *db.Database
	Logger *log.Logger
	Robot  *robot.Robot
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
		Robot:  robot.NewRobot(),
	}
}

func (s *Services) Close() {
	s.DB.Close()
}
