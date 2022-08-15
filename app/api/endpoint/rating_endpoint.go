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

type RatingEndpoint struct {
	CreateRatingTypeNum     endpoint.Endpoint
	UpdateRatingById        endpoint.Endpoint
	GetRatingTypeNumById    endpoint.Endpoint
	DeleteRatingTypeNumById endpoint.Endpoint
	GetRatingTypeNums       endpoint.Endpoint

	CreateRatingSubmission                  endpoint.Endpoint
	UpdateRatingSubmission                  endpoint.Endpoint
	GetRatingSubmission                     endpoint.Endpoint
	GetListRatingSubmission                 endpoint.Endpoint
	DeleteRatingSubmission                  endpoint.Endpoint
	GetListRatingSubmissionWithUserIdLegacy endpoint.Endpoint
	UpdateRatingSubDisplayNameByIdLegacy    endpoint.Endpoint

	CreateRatingTypeLikert     endpoint.Endpoint
	GetRatingTypeLikertById    endpoint.Endpoint
	UpdateRatingTypeLikertById endpoint.Endpoint
	DeleteRatingTypeLikertById endpoint.Endpoint
	GetRatingTypeLikerts       endpoint.Endpoint

	CreateRating                  endpoint.Endpoint
	ShowRating                    endpoint.Endpoint
	UpdateRating                  endpoint.Endpoint
	DeleteRating                  endpoint.Endpoint
	GetRatings                    endpoint.Endpoint
	GetListRatingSummary          endpoint.Endpoint
	GetRatingBySourceTypeAndActor endpoint.Endpoint

	CreateRatingFormula     endpoint.Endpoint
	UpdateRatingFormulaById endpoint.Endpoint
	GetRatingFormulaById    endpoint.Endpoint
	DeleteRatingFormulaById endpoint.Endpoint
	GetRatingFormulas       endpoint.Endpoint

	CreateRatingSubHelpful endpoint.Endpoint
}

func MakeRatingEndpoints(s service.RatingService) RatingEndpoint {
	return RatingEndpoint{
		CreateRatingTypeNum:     makeCreateRatingTypeNum(s),
		UpdateRatingById:        makeUpdateRatingById(s),
		GetRatingTypeNumById:    makeGetRatingTypeNumeById(s),
		DeleteRatingTypeNumById: makeDeleteRatingTypeNumById(s),
		GetRatingTypeNums:       makeGetRatingTypeNums(s),

		CreateRatingSubmission:                  makeCreateRatingSubmission(s),
		UpdateRatingSubmission:                  makeUpdateRatingSubmission(s),
		GetRatingSubmission:                     makeGetRatingSubmission(s),
		GetListRatingSubmission:                 makeGetListRatingSubmissions(s),
		DeleteRatingSubmission:                  makeDeleteRatingSubmission(s),
		GetListRatingSubmissionWithUserIdLegacy: makeGetListRatingSubmissionWithUserIdLegacy(s),
		UpdateRatingSubDisplayNameByIdLegacy:    makeUpdatePublicRatingSubDisplayNameByIdLegacy(s),

		CreateRatingTypeLikert:     makeCreateRatingTypeLikert(s),
		GetRatingTypeLikertById:    makeGetRatingTypeLikertById(s),
		UpdateRatingTypeLikertById: makeUpdateRatingTypeLikertById(s),
		DeleteRatingTypeLikertById: makeDeleteRatingTypeLikertById(s),
		GetRatingTypeLikerts:       makeRatingTypeLikerts(s),

		CreateRating:                  makeCreateRating(s),
		ShowRating:                    makeShowRating(s),
		UpdateRating:                  makeUpdateRating(s),
		DeleteRating:                  makeDeleteRatingById(s),
		GetRatings:                    makeGetListRatings(s),
		GetListRatingSummary:          makGetListRatingSummary(s),
		GetRatingBySourceTypeAndActor: makeGetRatingBySourceTypeAndActor(s),

		CreateRatingFormula:     makeCreateRatingFormula(s),
		UpdateRatingFormulaById: makeUpdateRatingFormulaById(s),
		GetRatingFormulaById:    makeGetRatingFormulaById(s),
		DeleteRatingFormulaById: makeDeleteRatingFormulaById(s),
		GetRatingFormulas:       makeRatingFormulas(s),

		CreateRatingSubHelpful: makeCreateRatingSubHelpful(s),
	}
}

func makeCreateRatingTypeNum(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.CreateRatingTypeNumRequest)

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

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

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

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

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

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

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

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

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		result, pagination, msg := s.GetRatingTypeNums(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeCreateRatingSubmission(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.CreateRatingSubmissionRequest)

		jwtObj, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		// Validate jwtObj User Id
		if jwtObj.UserIdLegacy != *req.UserIDLegacy {
			msg := message.ErrUserNotFound
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		// Set value from extract JWT
		if *req.DisplayName == "" {
			req.DisplayName = &jwtObj.Fullname
		}
		req.Avatar = jwtObj.Avatar

		result, msg := s.CreateRatingSubmission(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeUpdateRatingSubmission(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.UpdateRatingSubmissionRequest)

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

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

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
	}
}

func makeGetListRatingSubmissionWithUserIdLegacy(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetPublicListRatingSubmissionByUserIdRequest)

		jwtObj, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		// Validate jwtObj User Id
		if jwtObj.UserIdLegacy != req.UserIdLegacy {
			msg := message.ErrUserNotFound
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		result, pagination, msg := s.GetListRatingSubmissionWithUserIdLegacy(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, pagination), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeUpdatePublicRatingSubDisplayNameByIdLegacy(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.UpdateRatingSubDisplayNameRequest)

		jwtObj, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}
		// Validate jwtObj User Id
		if jwtObj.UserIdLegacy != req.UserIdLegacy {
			msg := message.ErrUserNotFound
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		msg := s.UpdateRatingSubDisplayNameByIdLegacy(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}

func makeGetRatingSubmission(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		result, msg := s.GetRatingSubmission(fmt.Sprint(rqst))

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

func makeGetListRatingSubmissions(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.ListRatingSubmissionRequest)

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

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

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

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

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

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

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

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

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

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

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

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

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		_, msg := s.CreateRating(req)

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}

func makeShowRating(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		result, msg := s.GetRatingById(fmt.Sprint(rqst))
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeUpdateRating(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		req := rqst.(request.UpdateRatingRequest)
		msg := s.UpdateRating(req)

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}

func makeDeleteRatingById(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		msg := s.DeleteRating(fmt.Sprint(rqst))
		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}

func makeGetListRatings(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetListRatingsRequest)

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

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

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		result, msg := s.GetListRatingSummary(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeGetRatingBySourceTypeAndActor(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetRatingBySourceTypeAndActorRequest)

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		result, msg := s.GetRatingBySourceTypeAndActor(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeCreateRatingFormula(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveRatingFormula)

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		result, msg := s.CreateRatingFormula(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeUpdateRatingFormulaById(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveRatingFormula)

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		msg := s.UpdateRatingFormula(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
	}
}

func makeGetRatingFormulaById(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetRatingFormulaRequest)

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		result, msg := s.GetRatingFormulaById(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeDeleteRatingFormulaById(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetRatingFormulaRequest)

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		msg := s.DeleteRatingFormulaById(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
	}
}

func makeRatingFormulas(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetRatingFormulasRequest)

		_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		result, pagination, msg := s.GetRatingFormulas(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeCreateRatingSubHelpful(s service.RatingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.CreateRatingSubHelpfulRequest)

		jwtObj, jwtMsg := global.SetJWTInfoFromContext(ctx)
		if jwtMsg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		}

		// Validate jwtObj User Id
		if jwtObj.UserIdLegacy != req.UserIDLegacy {
			msg := message.ErrUserNotFound
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		result, msg := s.CreateRatingSubHelpful(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}
