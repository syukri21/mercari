package grpc

import (
	"context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

// GrpcServer ...
type GrpcServer struct {
	getAreaInfo grpctransport.Handler
}

func (g GrpcServer) GetAreaInfo(ctx context.Context, request *GetAreaRequest) (*GetAreaInfoResponse, error) {
	_, resp, err := g.getAreaInfo.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*GetAreaInfoResponse), nil
}

func (g GrpcServer) mustEmbedUnimplementedServiceAreaServer() {
	//TODO implement me
	panic("implement me")
}
