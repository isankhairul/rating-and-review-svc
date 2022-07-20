package endpoint

import (
	"context"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"

	"github.com/go-kit/kit/endpoint"
)

type PublicRatingEndpoint struct {
	GetListRatingSubmissionBySourceTypeAndUID endpoint.Endpoint
	GetListRatingSummaryBySourceType          endpoint.Endpoint
}

func MakePublicRatingEndpoints(s service.PublicRatingService) PublicRatingEndpoint {
	return PublicRatingEndpoint{
		GetListRatingSummaryBySourceType:          makeGetListRatingSummaryBySourceType(s),
		GetListRatingSubmissionBySourceTypeAndUID: makeGetListRatingSubmissionBySourceTypeAndUID(s),
	}
}

func makeGetListRatingSummaryBySourceType(s service.PublicRatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetPublicListRatingSummaryRequest)
		result, pagination, msg := s.GetListRatingSummaryBySourceType(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, pagination), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeGetListRatingSubmissionBySourceTypeAndUID(s service.PublicRatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetPublicListRatingSubmissionRequest)
		result, pagination, msg := s.GetListRatingSubmissionBySourceTypeAndUID(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, pagination), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}
