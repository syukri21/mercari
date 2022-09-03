package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	gw "github.com/syukri21/mercari/service_apigw/internal/gateway/service_area"
)

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClient interface
var ClientRegisterer = registerer("grpc-gateway-service-auth")

type registerer string

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
),
) {
	f(string(r), r.registerClients)
}

func (r registerer) registerClients(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
	// check the passed configuration and initialize the plugin
	name, ok := extra["name"].(string)
	if !ok {
		return nil, errors.New("wrong config")
	}
	if name != string(r) {
		return nil, fmt.Errorf("unknown register %s", name)
	}
	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http client
	cfg := parse(extra)
	if cfg == nil {
		return nil, errors.New("wrong config")
	}
	if cfg.name != string(r) {
		return nil, fmt.Errorf("unknown register %s", cfg.name)
	}
	return gw.New(ctx, cfg.ServiceAuth)
}

func init() {
	fmt.Println("Service area client plugin loaded.")
}

func main() {}

func parse(extra map[string]interface{}) *opts {
	name, ok := extra["name"].(string)
	if !ok {
		return nil
	}

	rawEs, ok := extra["endpoints"]
	if !ok {
		return nil
	}
	es, ok := rawEs.([]interface{})
	if !ok || len(es) < 1 {
		return nil
	}
	endpoints := make([]string, len(es))
	for i, e := range es {
		endpoints[i] = e.(string)
	}

	return &opts{
		name:        name,
		ServiceAuth: endpoints[0],
	}
}

type opts struct {
	name        string
	ServiceAuth string
}
