package main

import (
	"context"
	"net/http"
	"time"

	"github.com/renniemaharaj/project-list-go/internal/repository"
	"github.com/renniemaharaj/project-list-go/internal/router"

	"github.com/renniemaharaj/conveyor/pkg/conveyor"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
)

func main() {
	l := logger.New().Prefix("Backend")
	m := conveyor.CreateManager().Start() // todo: utilize conveyor elsewhere
	m.Start()

	// allow time for postgres to intialize
	time.Sleep(5 * time.Second)

	repos, err := repository.Get()
	if err != nil {
		l.Fatal(err)

	}

	// will automatically initalize tables
	if err := repository.InitializeDatabaseTables(context.Background(), repos); err != nil {
		panic(err)
	}

	demoData := true
	if demoData {
		repo, err := repository.Get()
		if err != nil {
			logger.New().Fatal(err)
		}
		err = repo.InsertSeededDemoData(context.Background())
		if err != nil {
			logger.New().Fatal(err)
		}
	}
	// setup chi router and start server
	r := router.SetupRouter()
	go func() {
		l.Info("Starting server on :8081")
		if err := http.ListenAndServe(":8081", r); err != nil {
			l.Fatal(err)
		}
	}()

	select {}
}
