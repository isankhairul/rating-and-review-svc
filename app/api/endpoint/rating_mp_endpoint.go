package endpoint

import (
	"context"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"

	"github.com/go-kit/kit/endpoint"
)

type RatingMpEndpoint struct {
	CreateRatingSubmission           endpoint.Endpoint
	GetRatingSubmission              endpoint.Endpoint
	GetListRatingSubmission          endpoint.Endpoint
	GetListRatingSummaryBySourceType endpoint.Endpoint
}

func MakeRatingMpEndpoints(s service.RatingMpService) RatingMpEndpoint {
	return RatingMpEndpoint{
		CreateRatingSubmission:           makeCreateRatingSubmissionMp(s),
		GetRatingSubmission:              makeGetRatingSubmissionMp(s),
		GetListRatingSubmission:          makeGetListRatingSubmissionsMp(s),
		GetListRatingSummaryBySourceType: makeGetListRatingSummaryMpBySourceType(s),
	}
}

func makeGetListRatingSummaryMpBySourceType(s service.RatingMpService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetListRatingSummaryRequest)
		result, pagination, msg := s.GetListRatingSummaryBySourceType(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, pagination), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeCreateRatingSubmissionMp(s service.RatingMpService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.CreateRatingSubmissionRequest)

		jwtObj, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		// Validate jwtObj User Id
		if req.UserIDLegacy != nil && jwtObj.UserIdLegacy != *req.UserIDLegacy {
			msg := message.ErrUserNotFound
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		// Set value from extract JWT
		if req.DisplayName == nil || *req.DisplayName == "" {
			name := jwtObj.Fullname.(string)
			req.DisplayName = &name
		}
		req.Avatar = jwtObj.Avatar.(string)

		// set user_id_legacy from token jwt
		userIdLegacy := fmt.Sprintf("%v", jwtObj.UserIdLegacy)
		req.UserIDLegacy = &userIdLegacy
		req.UserID = &userIdLegacy

		result, msg := s.CreateRatingSubmissionMp(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeGetRatingSubmissionMp(s service.RatingMpService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		result, msg := s.GetRatingSubmissionMp(fmt.Sprint(rqst))

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeGetListRatingSubmissionsMp(s service.RatingMpService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.ListRatingSubmissionRequest)

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		result, pagination, msg := s.GetListRatingSubmissionsMp(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}
