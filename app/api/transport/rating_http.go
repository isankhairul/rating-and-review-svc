package transport

import (
	"context"
	"encoding/json"
	"github.com/gorilla/schema"
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

func RatingHttpHandler(s service.RatingService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeRatingEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
	}

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + "rating-types-numeric/").Handler(httptransport.NewServer(
		ep.CreateRatingTypeNum,
		decodeCreateRatingTypeNum,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "rating-types-numeric/{id}").Handler(httptransport.NewServer(
		ep.GetRatingTypeNumById,
		decodeGetRatingById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPut).Path(_struct.PrefixBase + "rating-types-numeric/{id}").Handler(httptransport.NewServer(
		ep.UpdateRatingById,
		decodeUpdateRatingById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodDelete).Path(_struct.PrefixBase + "rating-types-numeric/{id}").Handler(httptransport.NewServer(
		ep.DeleteRatingTypeNumById,
		decodeGetRatingById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "rating-types-numeric").Handler(httptransport.NewServer(
		ep.GetRatingTypeNums,
		decodeGetRatingTypeNums,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + "rating-submissions/").Handler(httptransport.NewServer(
		ep.CreateRatingSubmission,
		decodeCreateRatingSubmission,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "rating-submissions").Handler(httptransport.NewServer(
		ep.GetListRatingSubmission,
		decodeGetListRatingSubmission,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPut).Path(_struct.PrefixBase + "rating-submissions/{id}").Handler(httptransport.NewServer(
		ep.UpdateRatingSubmission,
		decodeUpdateRatingSubmission,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "rating-submissions/{id}").Handler(httptransport.NewServer(
		ep.GetRatingSubmission,
		decodeGetById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodDelete).Path(_struct.PrefixBase + "rating-submissions/{id}").Handler(httptransport.NewServer(
		ep.DeleteRatingSubmission,
		decodeGetById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + "rating-types-likert/").Handler(httptransport.NewServer(
		ep.CreateRatingTypeLikert,
		decodeCreateRatingTypeLikert,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "rating-types-likert/{id}").Handler(httptransport.NewServer(
		ep.GetRatingTypeLikertById,
		decodeGetRatingTypeLikertById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPut).Path(_struct.PrefixBase + "rating-types-likert/{id}").Handler(httptransport.NewServer(
		ep.UpdateRatingTypeLikertById,
		decodeUpdateRatingTypeLikertById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodDelete).Path(_struct.PrefixBase + "rating-types-likert/{id}").Handler(httptransport.NewServer(
		ep.DeleteRatingTypeLikertById,
		decodeGetRatingTypeLikertById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "rating-types-likert").Handler(httptransport.NewServer(
		ep.GetRatingTypeLikerts,
		decodeRatingTypeLikerts,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + _struct.PrefixRating + "/").Handler(httptransport.NewServer(
		ep.CreateRating,
		decodeCreateRating,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + "ratings/summary").Handler(httptransport.NewServer(
		ep.GetListRatingSummary,
		decodeGetRatingSummary,
		encoder.EncodeResponseHTTP,
		options...,
	))

	ratingPathUid := _struct.PrefixBase + _struct.PrefixRating + "/{id}"
	pr.Methods(http.MethodGet).Path(ratingPathUid).Handler(httptransport.NewServer(
		ep.ShowRating,
		decodeGetById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPut).Path(ratingPathUid).Handler(httptransport.NewServer(
		ep.UpdateRating,
		decodeEditRatingById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodDelete).Path(ratingPathUid).Handler(httptransport.NewServer(
		ep.DeleteRating,
		decodeGetById,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(_struct.PrefixBase + _struct.PrefixRating).Handler(httptransport.NewServer(
		ep.GetRatings,
		decodeGetRatings,
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
	err = req.Validate()
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
	var req request.CreateRatingSubmissonRequest
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
	var req request.UpdateRatingSubmissonRequest
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
	err = req.Validate()
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

	return req, nil
}
