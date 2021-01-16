package port

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	pb "github.com/zale144/portpb"
	"github.com/zale144/ports/clientapi/internal/service/port/mocks"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/zale144/ports/clientapi/pkg/apierror"
)

func TestProcessor_Process(t *testing.T) {
	type fields struct {
		client    client
		batchSize int64
	}
	type args struct {
		ctx  context.Context
		file io.Reader
	}

	ctx := context.Background()

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErrRet apierror.ErrorMessage
	}{
		{
			name:       "success",
			fields:     fields{
				client:    func () client {
					c := &mocks.MockClient{}
					c.On("ProcessPortBatch", ctx,
						mock.AnythingOfType("*proto.PortBatchReq")).
						Return(nil, nil)
					return c
				}(),
				batchSize: 2,
			},
			args:       args{
				ctx:  ctx,
				file: strings.NewReader(jsonFile),
			},
			wantErrRet: apierror.ErrorMessage{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Service{
				client:    tt.fields.client,
				batchSize: tt.fields.batchSize,
			}

			gotErrRet := u.Process(tt.args.ctx, tt.args.file)
			require.Equal(t, tt.wantErrRet.R, gotErrRet.R)
		})
	}
}

func TestService_GetPorts(t *testing.T) {
	type fields struct {
		client    client
		batchSize int64
	}
	type args struct {
		ctx context.Context
	}

	ctx := context.Background()
	var ports []*pb.Port
	_ = json.Unmarshal([]byte(jsonFile), &ports)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*pb.Port
		want1  apierror.ErrorMessage
	}{
		{
			name: "success",
			fields: fields{
				client:    func () client {
					c := &mocks.MockClient{}
					c.On("GetPorts", ctx,
						mock.AnythingOfType("*proto.GetPortsReq")).
						Return(&pb.GetPortsRsp{
						Ports:                ports,
					}, nil)
					return c
				}(),
				batchSize: 0,
			},
			args: args{
				ctx: ctx,
			},
			want: nil,
			want1: apierror.ErrorMessage{
				Message: "",
				R: apierror.List{
					Errors: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Service{
				client:    tt.fields.client,
				batchSize: tt.fields.batchSize,
			}
			got, got1 := p.GetPorts(tt.args.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPorts() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetPorts() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

const jsonFile = `{
  "AEAJM": {
    "name": "Ajman",
    "city": "Ajman",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "coordinates": [
      55.5136433,
      25.4052165
    ],
    "province": "Ajman",
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAJM"
    ],
    "code": "52000"
  },
  "AEAUH": {
    "name": "Abu Dhabi",
    "coordinates": [
      54.37,
      24.47
    ],
    "city": "Abu Dhabi",
    "province": "Abu Z¸aby [Abu Dhabi]",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAUH"
    ],
    "code": "52001"
  },
  "AEDXB": {
    "name": "Dubai",
    "coordinates": [
      55.27,
      25.25
    ],
    "city": "Dubai",
    "province": "Dubayy [Dubai]",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEDXB"
    ],
    "code": "52005"
  },
  "AEFJR": {
    "name": "Al Fujayrah",
    "coordinates": [
      56.33,
      25.12
    ],
    "city": "Al Fujayrah",
    "province": "Al Fujayrah",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEFJR"
    ]
  },
  "AEJEA": {
    "name": "Jebel Ali",
    "city": "Jebel Ali",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "coordinates": [
      55.0272904,
      24.9857145
    ],
    "province": "Dubai",
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEJEA"
    ],
    "code": "52051"
  }
}`
