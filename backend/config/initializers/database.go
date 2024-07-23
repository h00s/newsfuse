package initializers

import (
	"github.com/go-raptor/connector/postgres"
	"github.com/go-raptor/raptor/v2"
	"github.com/h00s/newsfuse/db"
)

func Database() raptor.Database {
	return raptor.Database{
		Connector:  postgres.NewPostgresConnector(),
		Migrations: db.Migrations(),
	}
}
