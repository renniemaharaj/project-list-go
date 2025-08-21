package main

import (
	"context"
	"net/http"

	"github.com/renniemaharaj/project-list-go/internal/repository"
	"github.com/renniemaharaj/project-list-go/internal/router"

	"github.com/renniemaharaj/conveyor/pkg/conveyor"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
)

func main() {
	r := router.SetupRouter()
	l := logger.New().Prefix("Backend")
	m := conveyor.CreateManager().Start()
	m.Start()

	if err := repository.InitDatabaseTables(context.Background()); err != nil {
		panic(err)
	}

	go func() {
		l.Info("Starting server on :8081")
		if err := http.ListenAndServe(":8081", r); err != nil {
			l.Fatal(err)
		}
	}()

	select {}
}
