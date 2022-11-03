package transport

import (
	"context"
	"encoding/json"
	"github.com/gorilla/schema"
	"go-klikdokter/app/api/endpoint"
	"go-klikdokter/app/model/base/encoder"
	request_dapr "go-klikdokter/app/model/request/dapr"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/_struct"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func DaprHttpHandler(s service.DaprService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeDaprEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + "/dapr/publisher").Handler(httptransport.NewServer(
		ep.Publisher,
		decodeDaprPublisher,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + "/dapr/subscriber/ratingsubcol").Handler(httptransport.NewServer(
		ep.SubscriberRatingsubcol,
		decodeDaprSubscriberRatingsubcol,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeDaprPublisher(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request_dapr.PublisherRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	return params, nil
}

func decodeDaprSubscriberRatingsubcol(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request_dapr.BodySubscriberRatingsubcolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}
