package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	gozipkin "github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	endpts "miaosha/sk-admin/endpoint"
	"miaosha/sk-admin/model"
	"net/http"
	"os"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

func MakeHttpHandler(ctx context.Context, endpoints endpts.SkAdminEndpoints, zipkinTracer *gozipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	zipkingServer := zipkin.HTTPServerTrace(zipkinTracer)

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkingServer,
	}

	r.Methods("Get").Path("/product/list").Handler(kithttp.NewServer(
		endpoints.GetProductEndpoint,
		decodeCreateProductCheckRequest,
		encodeResponse,
	))

	r.Methods("POST").Path("/product/create").Handler(kithttp.NewServer(
		endpoints.GetProductEndpoint,
		decodeCreateProductCheckRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/activity/create").Handler(kithttp.NewServer(
		endpoints.CreateActivityEndpoint,
		decodeCreateActivityCheckRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/activity/list").Handler(kithttp.NewServer(
		endpoints.GetActivityEndpoint,
		decodeGetListRequest,
		encodeResponse,
		options...,
	))

	r.Path("/metrics").Handler(promhttp.Handler())

	// create health check handler
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeResponse,
		options...,
	))

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	return loggedRouter
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// encodeArithmeticResponse encode response to return
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeCreateProductCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return nil, err
	}
	return product, nil
}

// decodeHealthCheckRequest decode request
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}

func decodeCreateActivityCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var activity model.Activity
	if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
		return nil, err
	}
	return activity, nil
}

// decodeUserRequest decode request params to struct
func decodeGetListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpts.GetListRequest{}, nil
}
