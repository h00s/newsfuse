package components

import (
	"github.com/go-raptor/connectors/bun/postgres"
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/db"
)

func New(c *raptor.Config) *raptor.Components {
	return &raptor.Components{
		DatabaseConnector: postgres.NewPostgresConnector(c.DatabaseConfig, db.Migrations()),
		Controllers:       Controllers(),
		Services:          Services(c),
		Middlewares:       Middlewares(),
	}
}
