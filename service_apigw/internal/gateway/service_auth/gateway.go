package service_area

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	serviceGrpc "github.com/syukri21/mercari/service_apigw/internal/generated/service_auth"
	"github.com/syukri21/mercari/service_apigw/internal/headermatcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// New http handler for every endpoint that are registered in template
func New(ctx context.Context, endpointPath string) (http.Handler, error) {
	gw := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
		runtime.WithIncomingHeaderMatcher(headermatcher.MatcherTwo),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			nrgrpc.UnaryClientInterceptor,
		),
		grpc.WithChainStreamInterceptor(
			nrgrpc.StreamClientInterceptor,
		),
	}

	if err := serviceGrpc.RegisterServiceAuthHandlerFromEndpoint(ctx, gw, endpointPath, opts); err != nil {
		return nil, err
	}

	return gw, nil
}
