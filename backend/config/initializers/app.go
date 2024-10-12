package initializers

import (
	"github.com/go-raptor/connector/bun/postgres"
	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/config"
	"github.com/h00s/newsfuse/db"
)

func App(c *raptor.Config) *raptor.AppInitializer {
	return &raptor.AppInitializer{
		Routes:            config.Routes(),
		DatabaseConnector: postgres.NewPostgresConnector(c.DatabaseConfig, db.Migrations()),
		Services:          Services(c),
		Middlewares:       Middlewares(),
		Controllers:       Controllers(),
	}
}
