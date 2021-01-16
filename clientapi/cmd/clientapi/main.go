package main

import (
	"github.com/zale144/ports/clientapi/internal/api/server"
	"github.com/zale144/ports/clientapi/internal/client"
	"github.com/zale144/ports/clientapi/internal/config"
	"github.com/zale144/ports/clientapi/internal/service/port"
	"log"
)

func main() {
	cfg, err := config.Configure()
	if err != nil {
		log.Fatal(err)
	}

	grpcConn, err := client.DialGrpc(cfg.GRPCURL)
	if err != nil {
		log.Fatal(err)
	}

	svc := port.NewService(grpcConn, cfg.BatchSize)
	server.ServeHTTP(cfg.HTTPPort, svc)
}
