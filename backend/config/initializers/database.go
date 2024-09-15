package initializers

import (
	"github.com/go-raptor/connector/bun/postgres"
	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/db"
)

func Database() raptor.Database {
	return raptor.Database{
		Connector:  postgres.NewPostgresConnector(),
		Migrations: db.Migrations(),
	}
}
