package publicendpoint

import (
	"context"
	"github.com/go-kit/log"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request/public"
	rp "go-klikdokter/app/repository"
	publicrepository "go-klikdokter/app/repository/public"
	"go-klikdokter/app/service/public"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"

	"github.com/go-kit/kit/endpoint"
)

type PublicRatingEndpoint struct {
	GetListRatingSubmissionBySourceTypeAndUID endpoint.Endpoint
	GetListRatingSummaryBySourceType          endpoint.Endpoint
	GetListDetailRatingSummaryBySourceType    endpoint.Endpoint
}

func MakePublicRatingEndpoints(s publicservice.PublicRatingService, logger log.Logger, db *mongo.Database) PublicRatingEndpoint {
	return PublicRatingEndpoint{
		GetListRatingSummaryBySourceType:          makeGetListRatingSummaryBySourceType(s, logger, db),
		GetListRatingSubmissionBySourceTypeAndUID: makeGetListRatingSubmissionBySourceTypeAndUID(s, logger, db),
		GetListDetailRatingSummaryBySourceType:    makeGetListDetailRatingSummaryBySourceType(s, logger, db),
	}
}

func makeGetListRatingSummaryBySourceType(s publicservice.PublicRatingService, logger log.Logger, db *mongo.Database) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(publicrequest.GetPublicListRatingSummaryRequest)
		var result interface{}
		var pagination *base.Pagination
		var msg message.Message

		if util.StringInSlice(strings.ToLower(req.SourceType), []string{"product", "store"}) {
			publicMp := publicservice.NewPublicRatingMpService(logger, rp.NewRatingMpRepository(db), publicrepository.NewPublicRatingMpRepository(db))
			result, pagination, msg = publicMp.GetListRatingSummaryBySourceType(req)
		} else {
			result, pagination, msg = s.GetListRatingSummaryBySourceType(req)
		}
		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, pagination), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeGetListRatingSubmissionBySourceTypeAndUID(s publicservice.PublicRatingService, logger log.Logger, db *mongo.Database) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(publicrequest.GetPublicListRatingSubmissionRequest)
		var result interface{}
		var pagination *base.Pagination
		var msg message.Message

		if util.StringInSlice(strings.ToLower(req.SourceType), []string{"product", "store"}) {
			publicMp := publicservice.NewPublicRatingMpService(logger, rp.NewRatingMpRepository(db), publicrepository.NewPublicRatingMpRepository(db))
			result, pagination, msg = publicMp.GetListRatingSubmissionBySourceTypeAndUID(req)
		} else {
			result, pagination, msg = s.GetListRatingSubmissionBySourceTypeAndUID(req)
		}

		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, pagination), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeGetListDetailRatingSummaryBySourceType(s publicservice.PublicRatingService, logger log.Logger, db *mongo.Database) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(publicrequest.PublicGetListDetailRatingSummaryRequest)
		var result interface{}
		var pagination *base.Pagination
		var msg message.Message

		publicMp := publicservice.NewPublicRatingMpService(logger, rp.NewRatingMpRepository(db), publicrepository.NewPublicRatingMpRepository(db))
		result, pagination, msg = publicMp.GetListDetailRatingSummaryBySourceType(req)

		if msg.Code != 212000 {
			return base.SetHttpResponse(msg.Code, msg.Message, encoder.Empty{}, pagination), nil
		}
		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}
