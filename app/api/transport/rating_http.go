package transport

import (
	"context"
	"encoding/json"
	"go-klikdokter/app/api/endpoint"
	"go-klikdokter/app/middleware"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request"
	publicrequest "go-klikdokter/app/model/request/public"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/_struct"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/schema"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func RatingHttpHandler(s service.RatingService, logger log.Logger, db *mongo.Database) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeRatingEndpoints(s, logger, db)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + "/rating-types-numeric/").Handler(httptransport.NewServer(
		ep.CreateRatingTypeNum,
		decodeCreateRatingTypeNum,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/rating-types-numeric/{id}").Handler(httptransport.NewServer(
		ep.GetRatingTypeNumById,
		decodeGetRatingById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPut).Path(_struct.PrefixBase + "/rating-types-numeric/{id}").Handler(httptransport.NewServer(
		ep.UpdateRatingById,
		decodeUpdateRatingById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodDelete).Path(_struct.PrefixBase + "/rating-types-numeric/{id}").Handler(httptransport.NewServer(
		ep.DeleteRatingTypeNumById,
		decodeGetRatingById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/rating-types-numeric").Handler(httptransport.NewServer(
		ep.GetRatingTypeNums,
		decodeGetRatingTypeNums,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + "/rating-submissions/").Handler(httptransport.NewServer(
		ep.CreateRatingSubmission,
		decodeCreateRatingSubmission,
		encoder.EncodeResponseHTTPWithCorrelationID,
		append(options, httptransport.ServerBefore(middleware.CorrelationIdToContext()))...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/rating-submissions").Handler(httptransport.NewServer(
		ep.GetListRatingSubmission,
		decodeGetListRatingSubmission,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPut).Path(_struct.PrefixBase + "/rating-submissions/{id}").Handler(httptransport.NewServer(
		ep.UpdateRatingSubmission,
		decodeUpdateRatingSubmission,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/rating-submissions/{id}").Handler(httptransport.NewServer(
		ep.GetRatingSubmission,
		decodeGetById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodDelete).Path(_struct.PrefixBase + "/rating-submissions/{id}").Handler(httptransport.NewServer(
		ep.DeleteRatingSubmission,
		decodeGetById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPut).Path(_struct.PrefixBase + "/cancel/rating-submissions").Handler(httptransport.NewServer(
		ep.CancelRatingSubByIds,
		decodeCancelRatingSub,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/list-rating-submissions/{source_type}/{source_uid}/{user_id_legacy}").Handler(httptransport.NewServer(
		ep.GetListRatingSubmissionWithUserIdLegacy,
		decodeGetRatingSubmissionWithUserIdLegacy,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPut).Path(_struct.PrefixBase + "/rating-submissions/user-id-legacy/{user_id_legacy}").Handler(httptransport.NewServer(
		ep.UpdateRatingSubDisplayNameByIdLegacy,
		decodeUpdatePublicRatingSubDisplayName,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPut).Path(_struct.PrefixBase + "/rating-submissions/reply/{id}").Handler(httptransport.NewServer(
		ep.ReplyAdminRatingSubmission,
		decodeReplyAdminRatingSubmission,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + "/rating-types-likert/").Handler(httptransport.NewServer(
		ep.CreateRatingTypeLikert,
		decodeCreateRatingTypeLikert,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/rating-types-likert/{id}").Handler(httptransport.NewServer(
		ep.GetRatingTypeLikertById,
		decodeGetRatingTypeLikertById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPut).Path(_struct.PrefixBase + "/rating-types-likert/{id}").Handler(httptransport.NewServer(
		ep.UpdateRatingTypeLikertById,
		decodeUpdateRatingTypeLikertById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodDelete).Path(_struct.PrefixBase + "/rating-types-likert/{id}").Handler(httptransport.NewServer(
		ep.DeleteRatingTypeLikertById,
		decodeGetRatingTypeLikertById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/rating-types-likert").Handler(httptransport.NewServer(
		ep.GetRatingTypeLikerts,
		decodeRatingTypeLikerts,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + "/ratings/").Handler(httptransport.NewServer(
		ep.CreateRating,
		decodeCreateRating,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/ratings/summary/{source_type}").Handler(httptransport.NewServer(
		ep.GetListRatingSummary,
		decodeGetRatingSummary,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/ratings/{id}").Handler(httptransport.NewServer(
		ep.ShowRating,
		decodeGetById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPut).Path(_struct.PrefixBase + "/ratings/{id}").Handler(httptransport.NewServer(
		ep.UpdateRating,
		decodeEditRatingById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodDelete).Path(_struct.PrefixBase + "/ratings/{id}").Handler(httptransport.NewServer(
		ep.DeleteRating,
		decodeGetById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/ratings").Handler(httptransport.NewServer(
		ep.GetRatings,
		decodeGetRatings,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/list-ratings/{source_type}/{source_uid}").Handler(httptransport.NewServer(
		ep.GetRatingBySourceTypeAndActor,
		decodeGetRatingBySourceTypeAndActor,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + "/rating-formula/").Handler(httptransport.NewServer(
		ep.CreateRatingFormula,
		decodeCreateRatingFormula,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/rating-formula").Handler(httptransport.NewServer(
		ep.GetRatingFormulas,
		decodeGetRatingFormulas,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "/rating-formula/{id}").Handler(httptransport.NewServer(
		ep.GetRatingFormulaById,
		decodeGetRatingFormulaById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPut).Path(_struct.PrefixBase + "/rating-formula/{id}").Handler(httptransport.NewServer(
		ep.UpdateRatingFormulaById,
		decodeUpdateRatingFormulaById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodDelete).Path(_struct.PrefixBase + "/rating-formula/{id}").Handler(httptransport.NewServer(
		ep.DeleteRatingFormulaById,
		decodeDeleteRatingFormulaById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + "/helpful-rating-submission/").Handler(httptransport.NewServer(
		ep.CreateRatingSubHelpful,
		decodeCreateRatingSubHelpful,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeCreateRatingTypeNum(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.CreateRatingTypeNumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	err = req.Validate()
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetRatingById(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.GetRatingTypeNumRequest
	req.Id = mux.Vars(r)["id"]
	return req, nil
}

func decodeUpdateRatingById(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.EditRatingTypeNumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	// err = req.Validate()
	if err != nil {
		return nil, err
	}
	req.Id = mux.Vars(r)["id"]
	return req, nil
}

func decodeGetRatingTypeNums(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.GetRatingTypeNumsRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	return params, nil
}

func decodeCreateRatingSubmission(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
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

func decodeUpdateRatingSubmission(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.UpdateRatingSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	err = req.Validate()

	if err != nil {
		return nil, err
	}

	req.ID = mux.Vars(r)["id"]

	return req, nil
}

func decodeGetById(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	uid := mux.Vars(r)["id"]
	return uid, nil
}

func decodeGetListRatingSubmission(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.ListRatingSubmissionRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	return params, nil
}

func decodeRatingTypeLikerts(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.GetRatingTypeLikertsRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	return params, nil
}

func decodeCreateRatingTypeLikert(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveRatingTypeLikertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	err = req.Validate()
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetRatingTypeLikertById(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.GetRatingTypeLikertRequest
	req.Id = mux.Vars(r)["id"]
	return req, nil
}

func decodeUpdateRatingTypeLikertById(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveRatingTypeLikertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	// err = req.Validate()
	if err != nil {
		return nil, err
	}
	req.Id = mux.Vars(r)["id"]
	return req, nil
}

func decodeCreateRating(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveRatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	err = req.Validate()
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeEditRatingById(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.UpdateRatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Body); err != nil {
		return nil, err
	}
	err = req.Body.Validate()
	if err != nil {
		return nil, err
	}
	req.Id = mux.Vars(r)["id"]
	return req, nil
}

func decodeGetRatings(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.GetListRatingsRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&req, r.Form); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetRatingSummary(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.GetListRatingSummaryRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&req, r.Form); err != nil {
		return nil, err
	}
	req.SourceType = mux.Vars(r)["source_type"]
	return req, nil
}

func decodeCreateRatingFormula(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveRatingFormula
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetRatingFormulas(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.GetRatingFormulasRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}
	return params, nil
}

func decodeGetRatingFormulaById(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.GetRatingFormulaRequest
	req.Id = mux.Vars(r)["id"]
	return req, nil
}

func decodeDeleteRatingFormulaById(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.GetRatingFormulaRequest
	req.Id = mux.Vars(r)["id"]
	return req, nil
}

func decodeUpdateRatingFormulaById(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveRatingFormula
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	req.Id = mux.Vars(r)["id"]
	return req, nil
}

func decodeGetRatingBySourceTypeAndActor(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req publicrequest.GetRatingBySourceTypeAndActorRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	if err = schema.NewDecoder().Decode(&req, r.Form); err != nil {
		return nil, err
	}

	req.SourceType = mux.Vars(r)["source_type"]
	req.SourceUID = mux.Vars(r)["source_uid"]

	return req, nil
}

func decodeGetRatingSubmissionWithUserIdLegacy(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.GetPublicListRatingSubmissionByUserIdRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}
	params.SourceType = mux.Vars(r)["source_type"]
	params.SourceUID = mux.Vars(r)["source_uid"]
	params.UserIdLegacy = mux.Vars(r)["user_id_legacy"]
	err = params.ValidateSourceType()
	if err != nil {
		return nil, err
	}
	return params, nil
}

func decodeUpdatePublicRatingSubDisplayName(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.UpdateRatingSubDisplayNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.UserIdLegacy = mux.Vars(r)["user_id_legacy"]
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

func decodeCancelRatingSub(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.CancelRatingById
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeReplyAdminRatingSubmission(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.ReplyAdminRatingSubmissionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.ID = mux.Vars(r)["id"]

	err = req.Validate()
	if err != nil {
		return nil, err
	}

	return req, nil
}
