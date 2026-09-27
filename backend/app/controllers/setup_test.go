package controllers_test

import (
	"os"
	"testing"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/config"
	"github.com/h00s/newsfuse/config/components"
)

var app *raptor.Raptor

func TestMain(m *testing.M) {
	app = raptor.NewTestApp(components.New(), config.Routes())
	os.Exit(m.Run())
}
