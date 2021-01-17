package port

import (
	"context"
	pb "github.com/zale144/portpb"
	"io"
	"log"

	"github.com/zale144/ports/portdomainservice/internal/model"
)

type Port struct {
	store store
}

//go:generate mockery --name=store --filename=store.go --structname=MockStore
type store interface {
	SavePorts(ctx context.Context, ports []model.Port) error
	GetPorts(ctx context.Context) ([]model.Port, error)
}

func NewPort(store store) Port {
	return Port{
		store: store,
	}
}

// ProcessPortBatch accepts a gRPC server stream with a list of port object to process and save to database
func (p Port) ProcessPortBatch(sv pb.PortService_ProcessPortBatchServer) error {

	ctx := sv.Context()

	for {
		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// receive data from stream
		in, err := sv.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return sv.SendAndClose(&pb.PortBatchRsp{Success: true})
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}

		firstID, lastID := in.Ports[0].Id, in.Ports[len(in.Ports)-1].Id
		log.Printf("received ports with id from %s to %s\n", firstID, lastID)

		// save the ports locally
		if err = p.store.SavePorts(ctx, toModel(in.Ports)); err != nil {
			log.Println(err)
			return err
		}
	}
}

func toModel(dtos []*pb.Port) []model.Port {
	var ports []model.Port
	for _, d := range dtos {
		ports = append(ports, model.Port{
			ID:          d.Id,
			Name:        d.Name,
			City:        d.City,
			Country:     d.Country,
			Alias:       d.Alias,
			Regions:     d.Regions,
			Coordinates: d.Coordinates,
			Province:    d.Province,
			Timezone:    d.Timezone,
			Unlocs:      d.Unlocs,
			Code:        d.Code,
		})
	}
	return ports
}

// GetPorts accepts a stream handle to push lists of ports to
func (p Port) GetPorts(_ *pb.GetPortsReq, sv pb.PortService_GetPortsServer) error {
	// TODO: paging
	ports, err := p.store.GetPorts(sv.Context())
	if err != nil {
		return err
	}
	err = sv.Send(&pb.GetPortsRsp{
		Ports: toPb(ports),
	})
	return nil
}

func toPb(ports []model.Port) []*pb.Port {
	var dtos []*pb.Port
	for _, d := range ports {
		dtos = append(dtos, &pb.Port{
			Id:          d.ID,
			Name:        d.Name,
			City:        d.City,
			Country:     d.Country,
			Alias:       d.Alias,
			Regions:     d.Regions,
			Coordinates: d.Coordinates,
			Province:    d.Province,
			Timezone:    d.Timezone,
			Unlocs:      d.Unlocs,
			Code:        d.Code,
		})
	}
	return dtos
}
