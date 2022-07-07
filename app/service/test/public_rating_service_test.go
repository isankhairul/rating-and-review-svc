package test

import (
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/message"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var publicRatingRepository = &repository_mock.PublicRatingRepositoryMock{Mock: mock.Mock{}}
var publicRactingService = service.NewPublicRatingService(logger, publicRatingRepository)

func init() {
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}
}

func TestGetRatingBySourceTypeAndSourceUID(t *testing.T) {
	req := request.GetRatingBySourceTypeAndActorRequest{
		SourceType: "doctor",
		SourceUID:  "894",
	}

	ratingId, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")
	ratingTypeId, _ := primitive.ObjectIDFromHex("62c4f03b6d90d90d6594fab5")
	description := "Description"
	statement01 := "Unsatisfied"
	statement02 := "Satisfied"
	resultRatings := []entity.RatingsCol{
		{
			ID:           ratingId,
			Name:         "Rating Kepuasan Layanan Dr. Yopie Tjandradiguna, sp.u",
			Description:  &description,
			SourceUid:    "894",
			SourceType:   "doctor",
			RatingTypeId: "62c4f03b6d90d90d6594fab5",
			RatingType:   "satisfied-unsatisfied-doctor",
		},
	}
	resultRatingType := entity.RatingTypesLikertCol{
		ID:            ratingTypeId,
		Type:          "satisfied-unsatisfied-doctor",
		Description:   &description,
		NumStatements: 2,
		Statement01:   &statement01,
		Statement02:   &statement02,
	}
	publicRatingRepository.Mock.On("GetRatingsBySourceTypeAndActor", req.SourceType, req.SourceUID).Return(resultRatings, nil).Once()
	publicRatingRepository.Mock.On("GetRatingTypeLikertById", ratingTypeId).Return(resultRatingType, nil).Once()

	_, msg := publicRactingService.GetRatingBySourceTypeAndActor(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetRatingBySourceTypeAndSourceUID_ErrNoDataRating(t *testing.T) {
	req := request.GetRatingBySourceTypeAndActorRequest{
		SourceType: "doctor",
		SourceUID:  "894",
	}
	resultRatings := []entity.RatingsCol{}
	publicRatingRepository.Mock.On("GetRatingsBySourceTypeAndActor", req.SourceType, req.SourceUID).Return(resultRatings, nil)

	_, msg := publicRactingService.GetRatingBySourceTypeAndActor(req)
	assert.Equal(t, message.ErrNoData, msg)
}
