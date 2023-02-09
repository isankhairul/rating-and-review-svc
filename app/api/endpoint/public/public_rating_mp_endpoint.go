package publicendpoint

import (
	"context"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/base/encoder"
	publicrequest "go-klikdokter/app/model/request/public"
	publicservice "go-klikdokter/app/service/public"

	"github.com/go-kit/kit/endpoint"
)

type PublicRatingMpEndpoint struct {
	GetListRatingSubmissionBySourceTypeAndUID endpoint.Endpoint
	GetListRatingSummaryBySourceType          endpoint.Endpoint
	GetListRatingSubmissionByID               endpoint.Endpoint
	GetRatingSummaryStoreProduct              endpoint.Endpoint
}

func MakePublicRatingMpEndpoints(s publicservice.PublicRatingMpService) PublicRatingMpEndpoint {
	return PublicRatingMpEndpoint{
		GetListRatingSummaryBySourceType:          makeGetListRatingSummaryMpBySourceType(s),
		GetListRatingSubmissionBySourceTypeAndUID: makeGetListRatingSubmissionMpBySourceTypeAndUID(s),
		GetListRatingSubmissionByID:               makeGetListRatingSubmissionByID(s),
		GetRatingSummaryStoreProduct:              makeGetRatingSummaryStoreProduct(s),
	}
}

func makeGetListRatingSummaryMpBySourceType(s publicservice.PublicRatingMpService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(publicrequest.GetPublicListRatingSummaryRequest)
		result, msg := s.GetListRatingSummaryBySourceType(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeGetListRatingSubmissionMpBySourceTypeAndUID(s publicservice.PublicRatingMpService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(publicrequest.GetPublicListRatingSubmissionRequest)
		result, pagination, msg := s.GetListRatingSubmissionBySourceTypeAndUID(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, pagination), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeGetListRatingSubmissionByID(s publicservice.PublicRatingMpService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(publicrequest.GetPublicListRatingSubmissionByIDRequest)
		result, pagination, msg, errMsg := s.GetListRatingSubmissionByID(ctx, req)
		if msg.Code != 212000 {
			return base.SetHttpResponseWithCorrelationID(ctx, msg.Code, msg.Message, nil, nil, errMsg), nil
		}
		return base.SetHttpResponseWithCorrelationID(ctx, msg.Code, msg.Message, result, pagination, errMsg), nil
	}
}

func makeGetRatingSummaryStoreProduct(s publicservice.PublicRatingMpService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(publicrequest.PublicGetRatingSummaryStoreProductRequest)
		result, msg := s.GetRatingSummaryStoreProduct(ctx, req)
		if msg.Code != 212000 {
			return base.SetHttpResponseWithCorrelationID(ctx, msg.Code, msg.Message, nil, nil, nil), nil
		}
		return base.SetHttpResponseWithCorrelationID(ctx, msg.Code, msg.Message, result, nil, nil), nil
	}
}
