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

	"github.com/gorilla/schema"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func RatingMpHttpHandler(s service.RatingMpService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeRatingMpEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/rating-submissions-mp").Handler(httptransport.NewServer(
		ep.GetListRatingSubmission,
		decodeGetListRatingSubmissionMp,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + "/rating-submissions-mp").Handler(httptransport.NewServer(
		ep.CreateRatingSubmission,
		decodeCreateRatingSubmissionMp,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/rating-submissions-mp/{id}").Handler(httptransport.NewServer(
		ep.GetRatingSubmission,
		decodeGetById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/ratings-summary-mp/{source_type}").Handler(httptransport.NewServer(
		ep.GetListRatingSummaryBySourceType,
		decodeGetRatingSummaryMpBySourceType,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeGetRatingSummaryMpBySourceType(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.GetListRatingSummaryRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}
	params.SourceType = mux.Vars(r)["source_type"]
	return params, nil
}

func decodeCreateRatingSubmissionMp(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.CreateRatingSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	err = req.Validate()

	if err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetListRatingSubmissionMp(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.ListRatingSubmissionRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	return params, nil
}
