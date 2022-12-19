package endpoint

import (
	"context"
	"go-klikdokter/app/model/base"
	request "go-klikdokter/app/model/request/upload"
	"go-klikdokter/app/service"

	"github.com/go-kit/log"

	"github.com/go-kit/kit/endpoint"
)

type UploadImageEndpoint struct {
	MakeUploadImage			endpoint.Endpoint
}

func MakeUploadImageEndpoint(s service.UploadService, logger log.Logger) UploadImageEndpoint {
	return UploadImageEndpoint{
		MakeUploadImage: 			makeUploadImage(s, logger),
	}
}

func makeUploadImage(s service.UploadService, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.UploadImageRequest)
		result, msg, errMsg := s.UploadImage(ctx, req)
		if msg.Code == 4000 {
			return base.SetHttpResponseWithCorrelationID(ctx, msg.Code, msg.Message, nil, nil, errMsg), nil
		}
		return base.SetHttpResponseWithCorrelationID(ctx, msg.Code, msg.Message, result, nil, errMsg), nil
	}
}