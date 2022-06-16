package endpoint

import (
	"context"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"

	"github.com/go-kit/kit/endpoint"
)

type RatingEndpoint struct {
	CreateRatingTypeNum     endpoint.Endpoint
	UpdateRatingById        endpoint.Endpoint
	GetRatingTypeNumById    endpoint.Endpoint
	DeleteRatingTypeNumById endpoint.Endpoint
	GetRatingTypeNums       endpoint.Endpoint

	CreateRatingSubmission  endpoint.Endpoint
	UpdateRatingSubmission  endpoint.Endpoint
	GetRatingSubmission     endpoint.Endpoint
	GetListRatingSubmission endpoint.Endpoint
	DeleteRatingSubmission  endpoint.Endpoint

	CreateRatingTypeLikert     endpoint.Endpoint
	GetRatingTypeLikertById    endpoint.Endpoint
	UpdateRatingTypeLikertById endpoint.Endpoint
	DeleteRatingTypeLikertById endpoint.Endpoint
	GetRatingTypeLikerts       endpoint.Endpoint

	CreateRating         endpoint.Endpoint
	ShowRating           endpoint.Endpoint
	UpdateRating         endpoint.Endpoint
	DeleteRating         endpoint.Endpoint
	GetRatings           endpoint.Endpoint
	GetListRatingSummary endpoint.Endpoint
}

func MakeRatingEndpoints(s service.RatingService) RatingEndpoint {
	return RatingEndpoint{
		CreateRatingTypeNum:     makeCreateRatingTypeNum(s),
		UpdateRatingById:        makeUpdateRatingById(s),
		GetRatingTypeNumById:    makeGetRatingTypeNumeById(s),
		DeleteRatingTypeNumById: makeDeleteRatingTypeNumById(s),
		GetRatingTypeNums:       makeGetRatingTypeNums(s),

		CreateRatingSubmission:  makeCreateRatingSubmission(s),
		UpdateRatingSubmission:  makeUpdateRatingSubmission(s),
		GetRatingSubmission:     makeGetRatingSubmission(s),
		GetListRatingSubmission: makeGetListRatingSubmissions(s),
		DeleteRatingSubmission:  makeDeleteRatingSubmission(s),

		CreateRatingTypeLikert:     makeCreateRatingTypeLikert(s),
		GetRatingTypeLikertById:    makeGetRatingTypeLikertById(s),
		UpdateRatingTypeLikertById: makeUpdateRatingTypeLikertById(s),
		DeleteRatingTypeLikertById: makeDeleteRatingTypeLikertById(s),
		GetRatingTypeLikerts:       makeRatingTypeLikerts(s),

		CreateRating:         makeCreateRating(s),
		ShowRating:           makeShowRating(s),
		UpdateRating:         makeUpdateRating(s),
		DeleteRating:         makeDeleteRatingById(s),
		GetRatings:           makeGetListRatings(s),
		GetListRatingSummary: makGetListRatingSummary(s),
	}
}

func makeCreateRatingTypeNum(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.CreateRatingTypeNumRequest)
		result, msg := s.CreateRatingTypeNum(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeUpdateRatingById(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.EditRatingTypeNumRequest)
		msg := s.UpdateRatingTypeNum(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
	}
}

func makeGetRatingTypeNumeById(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetRatingTypeNumRequest)
		result, msg := s.GetRatingTypeNumById(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeDeleteRatingTypeNumById(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetRatingTypeNumRequest)
		msg := s.DeleteRatingTypeNumById(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
	}
}

func makeGetRatingTypeNums(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetRatingTypeNumsRequest)
		result, pagination, msg := s.GetRatingTypeNums(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeCreateRatingSubmission(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.CreateRatingSubmissonRequest)
		msg := s.CreateRatingSubmission(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
	}
}

func makeUpdateRatingSubmission(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.UpdateRatingSubmissonRequest)
		msg := s.UpdateRatingSubmission(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
	}
}

func makeDeleteRatingSubmission(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		msg := s.DeleteRatingSubmission(fmt.Sprint(rqst))
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
	}
}

func makeGetRatingSubmission(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		result, msg := s.GetRatingSubmission(fmt.Sprint(rqst))
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeGetListRatingSubmissions(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.ListRatingSubmissionRequest)
		result, pagination, msg := s.GetListRatingSubmissions(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeCreateRatingTypeLikert(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveRatingTypeLikertRequest)
		msg := s.CreateRatingTypeLikert(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
	}
}

func makeGetRatingTypeLikertById(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetRatingTypeLikertRequest)
		result, msg := s.GetRatingTypeLikertById(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeUpdateRatingTypeLikertById(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveRatingTypeLikertRequest)
		msg := s.UpdateRatingTypeLikert(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
	}
}

func makeDeleteRatingTypeLikertById(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetRatingTypeLikertRequest)
		msg := s.DeleteRatingTypeLikertById(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
	}
}

func makeRatingTypeLikerts(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetRatingTypeLikertsRequest)
		result, pagination, msg := s.GetRatingTypeLikerts(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeCreateRating(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveRatingRequest)
		_, msg := s.CreateRating(req)

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}

func makeShowRating(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		result, msg := s.GetRatingById(fmt.Sprint(rqst))
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeUpdateRating(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.UpdateRatingRequest)
		msg := s.UpdateRating(req)

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}

func makeDeleteRatingById(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		msg := s.DeleteRating(fmt.Sprint(rqst))
		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}

func makeGetListRatings(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetListRatingsRequest)
		result, paging, msg := s.GetListRatings(req)

		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, paging), nil
	}
}

func makGetListRatingSummary(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetListRatingSummaryRequest)
		result, pagination, msg := s.GetListRatingSummary(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}