package publictransport

import (
	"context"
	publicendpoint "go-klikdokter/app/api/endpoint/public"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request/public"
	"go-klikdokter/app/service/public"
	"go-klikdokter/helper/_struct"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func PublicRatingHttpHandler(s publicservice.PublicRatingService, logger log.Logger, db *mongo.Database) http.Handler {
	pr := mux.NewRouter()

	ep := publicendpoint.MakePublicRatingEndpoints(s, logger, db)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
	}

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/public/ratings-summary/{source_type}").Handler(httptransport.NewServer(
		ep.GetListRatingSummaryBySourceType,
		decodeGetRatingSummaryBySourceType,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/public/rating-submissions/{source_type}/{source_uid}").Handler(httptransport.NewServer(
		ep.GetListRatingSubmissionBySourceTypeAndUID,
		decodeGetRatingSubmissionBySourceTypeAndUID,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeGetRatingSummaryBySourceType(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
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

func decodeGetRatingSubmissionBySourceTypeAndUID(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
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
