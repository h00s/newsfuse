package services

import "testing"

func TestDatabaseQueriesStopWithTheApp(t *testing.T) {
	res, shutdown := testResources()
	s := &DatabaseService{}
	if err := s.Init(res); err != nil {
		t.Fatal(err)
	}
	if err := s.Setup(); err != nil {
		t.Fatal(err)
	}

	shutdown()

	if s.Ctx.Err() == nil {
		t.Error("DB.Ctx is still live after shutdown, so a stray query could hold the pool open")
	}
}
