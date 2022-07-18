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
	GetRatingBySourceTypeAndActor             endpoint.Endpoint
	CreateRatingSubHelpful                    endpoint.Endpoint
	GetListRatingSummaryBySourceType          endpoint.Endpoint
	GetListRatingSubmissionBySourceTypeAndUID endpoint.Endpoint
	CreateRatingSubmission                    endpoint.Endpoint
	UpdateRatingSubDisplayNameByIdLegacy      endpoint.Endpoint
}

func MakePublicRatingEndpoints(s service.PublicRatingService) PublicRatingEndpoint {
	return PublicRatingEndpoint{
		GetRatingBySourceTypeAndActor:             makeGetRatingBySourceTypeAndActor(s),
		CreateRatingSubHelpful:                    makeCreateRatingSubHelpful(s),
		GetListRatingSummaryBySourceType:          makeGetListRatingSummaryBySourceType(s),
		GetListRatingSubmissionBySourceTypeAndUID: makeGetListRatingSubmissionBySourceTypeAndUID(s),
		CreateRatingSubmission:                    makeCreatePublicRatingSubmission(s),
		UpdateRatingSubDisplayNameByIdLegacy:      makeUpdatePublicRatingSubDisplayNameByIdLegacy(s),
	}
}

func makeGetRatingBySourceTypeAndActor(s service.PublicRatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetRatingBySourceTypeAndActorRequest)
		result, msg := s.GetRatingBySourceTypeAndActor(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeCreateRatingSubHelpful(s service.PublicRatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.CreateRatingSubHelpfulRequest)
		msg := s.CreateRatingSubHelpful(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
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

func makeCreatePublicRatingSubmission(s service.PublicRatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.CreateRatingSubmissionRequest)
		result, msg := s.CreatePublicRatingSubmission(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeUpdatePublicRatingSubDisplayNameByIdLegacy(s service.PublicRatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.UpdateRatingSubDisplayNameRequest)
		msg := s.UpdateRatingSubDisplayNameByIdLegacy(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}
