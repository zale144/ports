package main

import (
	"log"

	"github.com/zale144/ports/portdomainservice/internal/api/server"
	"github.com/zale144/ports/portdomainservice/internal/config"
	"github.com/zale144/ports/portdomainservice/internal/service/port"
	db "github.com/zale144/ports/portdomainservice/pkg/database"
	"github.com/zale144/ports/portdomainservice/repository"
)

func main() {

	cfg, err := config.Configure()
	if err != nil {
		log.Fatal(err)
	}

	dB, err := db.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPort(cfg, dB)
	svc := port.NewPort(repo)

	if err := server.Start(cfg, svc); err != nil {
		log.Fatal(err)
	}
}
