package test

import (
	"errors"
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
var publicRactingService = service.NewPublicRatingService(logger, ratingRepository, publicRatingRepository)

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

func TestCreateRatingSubHelpfulSuccess(t *testing.T) {
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	objectRatingSubHelpfulId, _ := primitive.ObjectIDFromHex("62c6438c08d23eb8fe9834e8")
	userId := "2210"
	comment := "Comment Test"

	input := request.CreateRatingSubHelpfulRequest{
		RatingSubmissionID: "629dce7bf1f26275e0d84826",
		UserID:             "2210",
		UserIDLegacy:       "2310",
		IPAddress:          "138.199.24.50",
		UserAgent:          "Chrome/{Chrome Rev} Mobile Safari/{WebKit Rev}",
	}

	ratingSubmission := entity.RatingSubmisson{
		ID: objectId,
		RatingID: "62c4f30f6d90d90d6594fab9	",
		UserID:        &userId,
		UserIDLegacy:  &userId,
		Comment:       &comment,
		Value:         "85",
		IPAddress:     "138.199.20.50",
		UserAgent:     "Chrome/{Chrome Rev} Mobile Safari/{WebKit Rev}",
		SourceTransID: "",
		LikeCounter:   3,
	}

	ratingSubHelpful := entity.RatingSubHelpfulCol{
		ID:                 objectRatingSubHelpfulId,
		RatingSubmissionID: "62c53baf039c7a6554accb0d",
		UserID:             "2210",
		UserIDLegacy:       "2310",
		IPAddress:          "138.199.20.50",
		UserAgent:          "Chrome/{Chrome Rev} Mobile Safari/{WebKit Rev}",
	}

	ratingRepository.Mock.On("GetRatingSubmissionById", objectId).Return(ratingSubmission, nil).Once()
	publicRatingRepository.Mock.On("CreateRatingSubHelpful", input).Return(ratingSubHelpful, nil).Once()
	publicRatingRepository.Mock.On("UpdateCounterRatingSubmission", objectId, ratingSubmission.LikeCounter).Return(nil).Once()

	msg := publicRactingService.CreateRatingSubHelpful(input)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestCreateRatingSubHelpfulRatingSubmissionNil(t *testing.T) {
	objectRatingSubmissionId, _ := primitive.ObjectIDFromHex("62c53baf039c7a6554accb0d")
	userId := "2210"
	comment := "Comment Test"

	input := request.CreateRatingSubHelpfulRequest{
		RatingSubmissionID: "62c53baf039c7a6554accb0d",
		UserID:             "2210",
		UserIDLegacy:       "2310",
		IPAddress:          "138.199.24.50",
		UserAgent:          "Chrome/{Chrome Rev} Mobile Safari/{WebKit Rev}",
	}

	ratingSubmission := entity.RatingSubmisson{
		ID: objectRatingSubmissionId,
		RatingID: "62c4f30f6d90d90d6594fab9	",
		UserID:        &userId,
		UserIDLegacy:  &userId,
		Comment:       &comment,
		Value:         "85",
		IPAddress:     "138.199.20.50",
		UserAgent:     "Chrome/{Chrome Rev} Mobile Safari/{WebKit Rev}",
		SourceTransID: "",
		LikeCounter:   3,
	}

	ratingRepository.Mock.On("GetRatingSubmissionById", objectRatingSubmissionId).Return(ratingSubmission, nil).Once()

	msg := publicRactingService.CreateRatingSubHelpful(input)
	assert.Equal(t, message.ErrRatingTypeNotExist, msg)
}

func TestCreateRatingSubHelpfulUpdateCounterFailed(t *testing.T) {
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	objectRatingSubHelpfulId, _ := primitive.ObjectIDFromHex("62c6438c08d23eb8fe9834e8")
	userId := "2210"
	comment := "Comment Test"

	input := request.CreateRatingSubHelpfulRequest{
		RatingSubmissionID: "629dce7bf1f26275e0d84826",
		UserID:             "2210",
		UserIDLegacy:       "2310",
		IPAddress:          "138.199.24.50",
		UserAgent:          "Chrome/{Chrome Rev} Mobile Safari/{WebKit Rev}",
	}

	ratingSubmission := entity.RatingSubmisson{
		ID: objectId,
		RatingID: "62c4f30f6d90d90d6594fab9	",
		UserID:        &userId,
		UserIDLegacy:  &userId,
		Comment:       &comment,
		Value:         "85",
		IPAddress:     "138.199.20.50",
		UserAgent:     "Chrome/{Chrome Rev} Mobile Safari/{WebKit Rev}",
		SourceTransID: "",
		LikeCounter:   3,
	}

	ratingSubHelpful := entity.RatingSubHelpfulCol{
		ID:                 objectRatingSubHelpfulId,
		RatingSubmissionID: "62c53baf039c7a6554accb0d",
		UserID:             "2210",
		UserIDLegacy:       "2310",
		IPAddress:          "138.199.20.50",
		UserAgent:          "Chrome/{Chrome Rev} Mobile Safari/{WebKit Rev}",
	}

	ratingRepository.Mock.On("GetRatingSubmissionById", objectId).Return(ratingSubmission, nil).Once()
	publicRatingRepository.Mock.On("CreateRatingSubHelpful", input).Return(ratingSubHelpful, nil).Once()
	publicRatingRepository.Mock.On("UpdateCounterRatingSubmission", objectId, ratingSubmission.LikeCounter).Return(errors.New("error")).Once()

	msg := publicRactingService.CreateRatingSubHelpful(input)
	assert.Equal(t, message.SuccessMsg, msg)
}
