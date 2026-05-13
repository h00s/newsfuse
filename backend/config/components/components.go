package components

import (
	"github.com/go-raptor/connectors/bun/postgres"
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/db"
)

func New() *raptor.Components {
	return &raptor.Components{
		DatabaseConnector: postgres.NewPostgresConnector(db.MigrationsFS()),
		Controllers:       Controllers(),
		Services:          Services(),
		Middlewares:       Middlewares(),
	}
}
