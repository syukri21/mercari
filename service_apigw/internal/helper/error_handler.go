// Package errhandler is a custom implementation of grpc-gateway error handler.
package errhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const msgGRPCConnIssue = "We're encountering a problem in our side, please try again later."

func setBaseHeaders(h http.Header, marshaler runtime.Marshaler) {
	h.Del("Trailer")
	h.Del("Transfer-Encoding")
	h.Set("Content-type", marshaler.ContentType("application/json"))
}

// Handler is a Stockbit's own grpc-gateway error handler.
func Handler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	setBaseHeaders(w.Header(), marshaler)
	grpcStatus := status.Convert(err)
	grpcCode := grpcStatus.Code()
	if grpcCode == codes.Unauthenticated {
		w.Header().Set("WWW-Authenticate", grpcStatus.Message())
	}

	encoder := json.NewEncoder(w)

	logTransportError(ctx, grpcStatus, r)

	// when got "Service Unavailable" from service discovery or service mesh:
	if grpcCode == codes.Unavailable && (strings.HasPrefix(grpcStatus.Message(), "no healthy upstream") || strings.HasPrefix(grpcStatus.Message(), "upstream connect error")) {
		httpCode := http.StatusBadGateway
		w.WriteHeader(httpCode) // override
		errJ := encoder.Encode(map[string]interface{}{
			"ErrorType": http.StatusBadGateway,
			"Message":   msgGRPCConnIssue,
		})
		if errJ == nil {
			return // we're done here
		}
		runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
		return // we're done here
	}

	w.WriteHeader(runtime.HTTPStatusFromCode(grpcCode))

	// if Stockbit's error details was attached, use it
	for _, d := range grpcStatus.Details() {
		grpcErr, ok := d.(*error)
		if !ok {
			continue
		}

		errJ := encoder.Encode(grpcErr)
		if errJ == nil {
			return // we're done here
		}

		runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
		return // we're done here
	}

	jErr := encoder.Encode(&map[string]interface{}{
		"ErrorType": http.StatusBadGateway,
		"Message":   msgGRPCConnIssue,
	})
	if jErr == nil {
		return // we're done here
	}

	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, fmt.Errorf("error %s", http.StatusBadGateway)) // final fallback handler
}

func logTransportError(ctx context.Context, stat *status.Status, r *http.Request) {
	switch stat.Code() {
	case codes.InvalidArgument, codes.Unavailable,
		codes.Aborted, codes.Canceled, codes.Unimplemented,
		codes.FailedPrecondition, codes.Unknown,
		codes.DeadlineExceeded, codes.ResourceExhausted:
	}
}
