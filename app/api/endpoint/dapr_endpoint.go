package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/base/encoder"
	request_dapr "go-klikdokter/app/model/request/dapr"
	"go-klikdokter/app/service"
)

type DaprEndpoint struct {
	Publisher              endpoint.Endpoint
	SubscriberRatingsubcol endpoint.Endpoint
}

func MakeDaprEndpoints(s service.DaprService) DaprEndpoint {
	return DaprEndpoint{
		Publisher:              makeDaprPublisher(s),
		SubscriberRatingsubcol: makeDaprSubscriberRatingsubcol(s),
	}
}

func makeDaprPublisher(s service.DaprService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request_dapr.PublisherRequest)

		//_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		//if jwtMsg.Code != message.SuccessMsg.Code {
		//	return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		//}

		result, msg := s.Publisher(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeDaprSubscriberRatingsubcol(s service.DaprService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request_dapr.BodySubscriberRatingsubcolRequest)

		//_, jwtMsg := global.SetJWTInfoFromContext(ctx)
		//if jwtMsg.Code != message.SuccessMsg.Code {
		//	return base.SetHttpResponse(jwtMsg.Code, jwtMsg.Message, nil, nil), nil
		//}

		result, msg := s.SubscriberRatingsubcol(req)
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}
