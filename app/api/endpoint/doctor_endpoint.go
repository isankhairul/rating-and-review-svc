package endpoint

import (
	"context"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"

	"github.com/go-kit/kit/endpoint"
)

type DoctorEndpoint struct {
	SaveDoctor endpoint.Endpoint
	Show       endpoint.Endpoint
}

func MakeDoctorEndpoints(s service.DoctorService) DoctorEndpoint {
	return DoctorEndpoint{
		SaveDoctor: makeSaveDoctor(s),
		Show:       makeShowDoctor(s),
	}
}

func makeSaveDoctor(s service.DoctorService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveDoctorRequest)
		result, msg := s.CreateDoctor(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeShowDoctor(s service.DoctorService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		result, msg := s.GetDoctor(fmt.Sprint(rqst))
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}
