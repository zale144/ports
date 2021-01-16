package port

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	pb "github.com/zale144/portpb"
	"github.com/zale144/ports/portdomainservice/internal/model"
	"github.com/zale144/ports/portdomainservice/internal/service/port/mocks"
	"reflect"
	"testing"
)

func TestPort_ProcessPortBatch(t *testing.T) {
	type fields struct {
		store store
	}
	type args struct {
		ctx context.Context
		in  *pb.PortBatchReq
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRsp *pb.PortBatchRsp
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				store: func() store {
					c := &mocks.MockStore{}
					c.On("SavePorts", ctx,
						mock.AnythingOfType("[]model.Port")).
						Return(nil)
					return c
				}(),
			},
			args: args{
				ctx: ctx,
				in: &pb.PortBatchReq{
					Ports: []*pb.Port{
						{
							Id:          "AEAJM",
							Name:        "Ajman",
							City:        "Ajman",
							Country:     "United Arab Emirates",
							Alias:       []string{},
							Regions:     []string{},
							Coordinates: []float64{55.5136433, 25.4052165},
							Province:    "Ajman",
							Timezone:    "Asia/Dubai",
							Unlocs:      []string{"AEAJM"},
							Code:        "52000",
						},
					},
				},
			},
			wantRsp: &pb.PortBatchRsp{
				Success: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Port{
				store: tt.fields.store,
			}
			gotRsp, err := p.ProcessPortBatch(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessPortBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.wantRsp, gotRsp)
		})
	}
}

func TestPort_GetPorts(t *testing.T) {
	type fields struct {
		store store
	}
	type args struct {
		ctx context.Context
		in1 *pb.GetPortsReq
	}

	ctx := context.Background()
	ports := []model.Port{
		{
			ID:          "AEAJM",
			Name:        "Ajman",
			City:        "Ajman",
			Country:     "United Arab Emirates",
			Alias:       []string{},
			Regions:     []string{},
			Coordinates: []float64{55.5136433, 25.4052165},
			Province:    "Ajman",
			Timezone:    "Asia/Dubai",
			Unlocs:      []string{"AEAJM"},
			Code:        "52000",
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GetPortsRsp
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				store: func() store {
					c := &mocks.MockStore{}
					c.On("GetPorts", ctx).
						Return(ports, nil)
					return c
				}(),
			},
			args: args{
				ctx: ctx,
				in1: &pb.GetPortsReq{},
			},
			want: &pb.GetPortsRsp{
				Ports: toPb(ports),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Port{
				store: tt.fields.store,
			}
			got, err := p.GetPorts(tt.args.ctx, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPorts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPorts() got = %v, want %v", got, tt.want)
			}
		})
	}
}
