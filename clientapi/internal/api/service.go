package api

import (
	"context"
	pb "github.com/zale144/portpb"
	"io"

	"github.com/zale144/ports/clientapi/pkg/apierror"
)

type PortService interface {
	Process(ctx context.Context, file io.Reader) apierror.ErrorMessage
	GetPorts(ctx context.Context) ([]*pb.Port, apierror.ErrorMessage)
}
