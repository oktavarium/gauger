package grpcserver

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pbapi "github.com/oktavarium/go-gauger/api"
)

func (server *GrpcServer) GetAll(ctx context.Context, _ *emptypb.Empty) (*pbapi.GetAllResponse, error) {
	resp := &pbapi.GetAllResponse{}
	allMetrics, err := server.storage.GetAllMetrics(ctx)

	if err != nil {
		err = status.Errorf(codes.Internal, fmt.Sprintf("error on getting all metrics: %s", err))
		return resp, err
	}

	resp.Metrics = pbapi.ConvertDBMetricsToMetrics(allMetrics)
	return resp, nil
}
