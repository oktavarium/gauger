package grpcserver

import (
	"context"
	"fmt"

	pbapi "github.com/oktavarium/go-gauger/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *GrpcServer) Updates(ctx context.Context, in *pbapi.UpdatesRequest) (*emptypb.Empty, error) {
	allMetrics := pbapi.ConvertMetricsToDBMetrics(in.GetMetrics())
	if err := server.storage.BatchUpdate(ctx, allMetrics); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, fmt.Sprintf("error on batch update: %s", err))
	}
	return &emptypb.Empty{}, nil
}
