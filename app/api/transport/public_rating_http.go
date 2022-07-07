package transport

import (
	"context"
	"encoding/json"
	"go-klikdokter/app/api/endpoint"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/_struct"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func PublicRatingHttpHandler(s service.PublicRatingService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakePublicRatingEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
	}

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "public/ratings/{source_type}/{source_uid}").Handler(httptransport.NewServer(
		ep.GetRatingBySourceTypeAndActor,
		decodeGetRatingBySourceTypeAndActor,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + "public/helpful_rating_submission/").Handler(httptransport.NewServer(
		ep.CreateRatingSubHelpful,
		decodeCreateRatingSubHelpful,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeGetRatingBySourceTypeAndActor(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.GetRatingBySourceTypeAndActorRequest
	req.SourceType = mux.Vars(r)["source_type"]
	req.SourceUID = mux.Vars(r)["source_uid"]
	return req, nil
}

func decodeCreateRatingSubHelpful(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.CreateRatingSubHelpfulRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	err = req.Validate()
	if err != nil {
		return nil, err
	}
	return req, nil
}
