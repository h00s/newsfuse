package migrate

import (
	"context"

	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/app/models"
)

func AddHeadline(db *raptor.DB) error {
	_, err := db.NewCreateTable().
		Model((*models.Headline)(nil)).
		ForeignKey(`("source_id") REFERENCES "sources" ("id") ON DELETE CASCADE`).
		Exec(context.Background())
	return err

	// _, err := db.NewCreateIndex().
}
