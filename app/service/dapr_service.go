package service

import (
	"fmt"
	request_dapr "go-klikdokter/app/model/request/dapr"
	helper_dapr "go-klikdokter/helper/dapr"
	"go-klikdokter/helper/message"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type DaprService interface {
	// PUB/SUB
	Publisher(input request_dapr.PublisherRequest) (map[string]interface{}, message.Message)
	SubscriberRatingsubcol(input request_dapr.BodySubscriberRatingsubcolRequest) (string, message.Message)
}

type daprServiceImpl struct {
	logger log.Logger
}

func NewDaprService(
	lg log.Logger,
) DaprService {
	return &daprServiceImpl{lg}
}

// swagger:route POST /dapr/publisher DAPR PublisherRequest
// Publish Event
//
// security:
// - Bearer: []
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *daprServiceImpl) Publisher(input request_dapr.PublisherRequest) (map[string]interface{}, message.Message) {
	logger := log.With(s.logger, "DaprService", "Publisher")

	client := helper_dapr.NewDaprHttpClient()
	response, err := client.PublishEvent(input.Topic, input.Data)

	_ = level.Info(logger).Log(fmt.Sprintf("Publisher response: %v", response))
	_ = level.Info(logger).Log(fmt.Sprintf("Publisher error: %v", err))

	if err != nil {
		return response, message.ErrInternalError
	}

	return response, message.SuccessMsg
}

// swagger:route POST /dapr/subscriber/ratingsubcol DAPR SubscriberRatingsubcolRequest
// Subsriber for topic queuing.rnr.ratingsubcol
//
// security:
// - Bearer: []
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *daprServiceImpl) SubscriberRatingsubcol(input request_dapr.BodySubscriberRatingsubcolRequest) (string, message.Message) {
	logger := log.With(s.logger, "DaprService", "SubscriberRatingsubcol")

	_ = level.Info(logger).Log(fmt.Sprintf("Got message: %v", input.Data))

	return "Hello " + input.Data, message.SuccessMsg
}
