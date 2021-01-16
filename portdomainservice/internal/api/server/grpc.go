package server

import (
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"
	pb "github.com/zale144/portpb"
	"google.golang.org/grpc"

	"github.com/zale144/ports/portdomainservice/internal/config"
)

func Start(cfg *config.Config, svc pb.PortServiceServer) error {

	url := fmt.Sprintf(":%s", cfg.GRPCPort)
	lis, err := net.Listen("tcp", url)
	if err != nil {
		return errors.Wrap(err, "failed to listen gRPC")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPortServiceServer(grpcServer, svc)

	log.Printf("gRPC listening on: %s", url)
	if err := grpcServer.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to serve gRPC")
	}
	return nil
}
