package port

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	pb "github.com/zale144/portpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"math"
	"runtime"

	"github.com/zale144/ports/clientapi/pkg/apierror"
)

type Service struct {
	client    client
	batchSize int64
}

//go:generate mockery --name=client --filename=client.go --structname=MockClient
type client interface {
	ProcessPortBatch(ctx context.Context, opts ...grpc.CallOption) (pb.PortService_ProcessPortBatchClient, error)
	GetPorts(ctx context.Context, in *pb.GetPortsReq, opts ...grpc.CallOption) (pb.PortService_GetPortsClient, error)
}

func NewService(conn *grpc.ClientConn, bs int64) Service {
	return Service{
		client:    pb.NewPortServiceClient(conn),
		batchSize: bs,
	}
}
// Process takes json file, processes batches of port objects and sends each batch to portdomainservice
// for them to be saved to the database
func (p Service) Process(ctx context.Context, file io.Reader) (errRet apierror.ErrorMessage) {

	defer func() {
		if !errRet.R.HasErrors() {
			errRet.Message = "completed successfully"
		} else {
			errRet.Message = "completed with errors"
		}
	}()

	// number of workers from number of CPU cores - minus one for the main goroutine
	numWorkers := int(math.Max(1.0, float64(runtime.NumCPU()-1)))

	// prepare channels for communicating parsed data and termination
	batches, ports, done := make(chan []pb.Port, numWorkers), make(chan *pb.PortBatchReq, numWorkers), make(chan int, numWorkers)

	// start the number of workers determined by numWorkers
	fmt.Printf("Starting %v workers...\n", numWorkers)
	for i := 0; i < numWorkers; i++ {
		go prepareBatch(ctx, i, batches, ports, done)
	}

	go p.readJSON(ctx, file, &errRet, batches)

	waits := numWorkers

	for {
		select {
		case portBatch, ok := <-ports:
			if !ok {
				return
			}
			// send to port domain service via gRPC
			log.Printf("sening ports %v\n", portBatch.Ports)
			sv, err := p.client.ProcessPortBatch(ctx)
			if err != nil {
				log.Printf("error getting server stream %v\n", err)
				errRet.R.AddError(err)
			}

			if err = sv.Send(portBatch); err != nil {
				log.Printf("error sending batch to portdomainservice %v\n", err)
				errRet.R.AddError(err)
			} else {
				log.Printf("sent [%d] ports to portdomainservice\n", len(portBatch.Ports))
			}

		case <-done:
			waits--
			if waits == 0 {
				close(ports)
			}
		// expect cancel
		case <-ctx.Done():
			errRet.Message = "cancelled by caller"
			return
		}
	}
}

func (p Service) readJSON(ctx context.Context, file io.Reader, errRet *apierror.ErrorMessage, batches chan []pb.Port) {
	dec := json.NewDecoder(file)

	t, err := dec.Token()
	if err != nil {
		log.Println(err)
		errRet.R.AddError(err)
		return
	}
	delim, ok := t.(json.Delim)
	if !ok || delim != '{' {
		errRet.R.AddError(errors.New("expected object"))
		return
	}
	defer close(batches)

	batch := make([]pb.Port, 0)

	for dec.More() {
		// expect cancel
		select {
		case <-ctx.Done():
			return
		default:
			t, err = dec.Token()
			if err != nil {
				errRet.R.AddError(err)
				return
			}
			if t != "" {
				var port pb.Port
				err = dec.Decode(&port)
				if err != nil {
					errRet.R.AddError(err)
					return
				}
				port.Id = t.(string)
				log.Printf("port %v", port)
				batch = append(batch, port)
				// if batch is not full
				if dec.More() && len(batch) < int(p.batchSize) {
					// keep adding to it
					continue
				}
			}
			// otherwise, send the batch for further processing
			batches <- batch
			batch = nil
			continue
		}
	}
}

func prepareBatch(ctx context.Context, id int, batches <-chan []pb.Port, portsReq chan<- *pb.PortBatchReq, done chan<- int) {
	proc := 0 // how many batches did we process?

	for batch := range batches {
		ports := make([]*pb.Port, 0)
		for i := range batch {
			// expect cancel
			select {
			case <-ctx.Done():
				return
			default:
				ports = append(ports, &batch[i])
			}
		}
		portsReq <- &pb.PortBatchReq{
			Ports: ports,
		}

		proc++

	}
	log.Printf("worker %d finished after processing %d batches\n", id, proc)
	done <- id // send goroutine identifier to done channel
}
// GetPorts uses the grpc stream client to get lists of ports
func (p Service) GetPorts(ctx context.Context) ([]*pb.Port, apierror.ErrorMessage) {
	errRsp := apierror.ErrorMessage{}
	stream, err := p.client.GetPorts(ctx, &pb.GetPortsReq{})
	if err != nil {
		errRsp.R.AddError(err)
		return nil, errRsp
	}

	portsCh := make(chan []*pb.Port)

	go func(ch chan []*pb.Port) {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(ch)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			ch <- resp.Ports
		}
	}(portsCh)

	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			log.Println(err)
		}
	}()

	var ports []*pb.Port
	for p := range portsCh {
		ports = append(ports, p...)
	}
	return ports, errRsp
}
