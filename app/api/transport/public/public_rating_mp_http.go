package publictransport

import (
	"context"
	publicendpoint "go-klikdokter/app/api/endpoint/public"
	"go-klikdokter/app/model/base/encoder"
	publicrequest "go-klikdokter/app/model/request/public"
	publicservice "go-klikdokter/app/service/public"
	"go-klikdokter/helper/_struct"
	"net/http"

	"github.com/gorilla/schema"

	"go-klikdokter/app/middleware"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func PublicRatingMpHttpHandler(s publicservice.PublicRatingMpService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := publicendpoint.MakePublicRatingMpEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
	}
	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/public/ratings-summary-mp/{source_type}").Handler(httptransport.NewServer(
		ep.GetListRatingSummaryBySourceType,
		decodeGetRatingSummaryMpBySourceType,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/public/rating-submissions-mp/{source_type}/{source_uid}").Handler(httptransport.NewServer(
		ep.GetListRatingSubmissionBySourceTypeAndUID,
		decodeGetRatingSubmissionMpBySourceTypeAndUID,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/public/rating-submissions-by-id").Handler(httptransport.NewServer(
		ep.GetListRatingSubmissionByID,
		decodeGetRatingSubmissionByID,
		encoder.EncodeResponseHTTPWithCorrelationID,
		append(options, httptransport.ServerBefore(middleware.CorrelationIdToContext()))...,
	))

	return pr
}

func decodeGetRatingSummaryMpBySourceType(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params publicrequest.GetPublicListRatingSummaryRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}
	params.SourceType = mux.Vars(r)["source_type"]
	return params, nil
}

func decodeGetRatingSubmissionMpBySourceTypeAndUID(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params publicrequest.GetPublicListRatingSubmissionRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}
	params.SourceType = mux.Vars(r)["source_type"]
	params.SourceUID = mux.Vars(r)["source_uid"]
	err = params.ValidateSourceType()
	if err != nil {
		return nil, err
	}
	return params, nil
}

func decodeGetRatingSubmissionByID(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params publicrequest.GetPublicListRatingSubmissionByIDRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}
	return params, nil
}