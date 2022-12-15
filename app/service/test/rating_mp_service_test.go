package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/message"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"testing"
)

var ratingMpRepository = &repository_mock.RatingMpRepository{Mock: mock.Mock{}}
var ratingMpSvc = service.NewRatingMpService(logger, ratingMpRepository)

func TestCreateRatingSubmissionMpSuccess(t *testing.T) {
	likertID := "629dce7bf1f26275e0d84820"
	objectLikertId, _ := primitive.ObjectIDFromHex(likertID)
	objectId, _ := primitive.ObjectIDFromHex(id)
	value := "80"

	input := request.CreateRatingSubmissionRequest{
		UserID:        &id,
		UserIDLegacy:  &id,
		DisplayName:   &name,
		SourceTransID: "8888888||kjkjkbjgjjh",
		SourceUID:     "9YUHJHHJKH99OJKJKJKJKJKJK",
		RatingType:    "rating_for_product",
		Value:         value,
	}

	sub := entity.RatingSubmissionMp{
		ID:            objectId,
		UserID:        &id,
		UserIDLegacy:  &id,
		RatingID:      id,
		SourceTransID: "8888888||kjkjkbjgjjh||629dce7bf1f26275e0d84826||629dce7bf1f26275e0d84826",
		Value:         value,
	}
	arrSub := []entity.RatingSubmissionMp{sub}

	likert := entity.RatingTypesLikertCol{
		ID:            objectLikertId,
		Description:   &description,
		NumStatements: 1,
		Statement01:   &description,
	}

	saveReq := []request.SaveRatingSubmissionMp{
		{
			UserID:        &id,
			UserIDLegacy:  &id,
			RatingID:      id,
			DisplayName:   &name,
			SourceTransID: "8888888||kjkjkbjgjjh||629dce7bf1f26275e0d84826||629dce7bf1f26275e0d84826",
			SourceUID:     "9YUHJHHJKH99OJKJKJKJKJKJK",
			Value:         &value,
		},
	}

	rating := entity.RatingsMpCol{
		ID:             objectId,
		RatingTypeId:   id,
		CommentAllowed: &Bool,
		Status:         &Bool,
		Description:    &name,
	}

	ratingMpRepository.Mock.On("FindRatingBySourceUIDAndRatingType", input.SourceUID, input.RatingType).Return(&rating, nil)
	ratingMpRepository.Mock.On("FindRatingSubmissionByUserIDLegacyAndRatingID", &id, id, sub.SourceTransID).Return(nil, gorm.ErrRecordNotFound)
	ratingMpRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", &id, id, sub.SourceTransID).Return(nil, gorm.ErrRecordNotFound)
	ratingMpRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", objectLikertId).Return(likert, nil)
	ratingMpRepository.Mock.On("GetRatingSubmissionById", objectId).Return(&sub, nil)
	ratingMpRepository.Mock.On("GetRatingById", objectId).Return(&rating, nil)
	ratingMpRepository.Mock.On("CreateRatingSubmission", saveReq).Return(&arrSub, nil)

	_, msg := ratingMpSvc.CreateRatingSubmissionMp(input)

	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetRatingSubmissionMpSuccess(t *testing.T) {
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	getSub := entity.RatingSubmissionMp{
		ID:      objectId,
		Comment: &Desc,
	}

	ratingMpRepository.Mock.On("GetRatingSubmissionById", objectId).Return(&getSub, nil)

	_, msg := ratingMpSvc.GetRatingSubmissionMp(id)

	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetRatingSubmissionMpFail(t *testing.T) {
	failId := "629dce7bf1f26275e0d84827"
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84827")

	ratingMpRepository.Mock.On("GetRatingSubmissionById", objectId).Return(nil, mongo.ErrNoDocuments)

	_, msg := ratingMpSvc.GetRatingSubmissionMp(failId)

	assert.Equal(t, message.ErrRatingSubmissionNotFound, msg)
}
