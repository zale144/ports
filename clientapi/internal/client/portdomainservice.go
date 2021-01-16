package client

import (
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"time"
)

func DialGrpc(grpcURL string) (*grpc.ClientConn, error) {

	retriableErrors := []codes.Code{codes.DataLoss}
	retryTimeout := 50 * time.Millisecond

	unaryInterceptor := grpc_retry.UnaryClientInterceptor(
		grpc_retry.WithCodes(retriableErrors...),
		grpc_retry.WithMax(3),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(retryTimeout)),
	)

	grpcConn, err := grpc.Dial(grpcURL,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		return nil, err
	}

	log.Println("Dial gRPC at", grpcURL)
	return grpcConn, nil
}
