package test

import (
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/message"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var ratingMpRepository = &repository_mock.RatingMpRepository{Mock: mock.Mock{}}
var ratingMpSvc = service.NewRatingMpService(logger, ratingMpRepository)

func TestCreateRatingSubmissionMpSuccess(t *testing.T) {
	userId := "34343432"
	orderNumber := "888888"
	value := "4"
	SourceTransID := orderNumber + "||product||Frtgffggffgft123||34343432"
	input := request.CreateRatingSubmissionRequest{
		UserID:        &userId,
		UserIDLegacy:  &userId,
		DisplayName:   &name,
		SourceTransID: orderNumber,
		SourceUID:     "Frtgffggffgft123",
		RatingType:    "rating_for_product",
		Value:         "4",
		MediaPath: []request.MediaPathObj{
			
		},
	}

	objectID, _ := primitive.ObjectIDFromHex(id)
	sub := entity.RatingSubmissionMp{
		ID:            objectID,
		UserID:        &userId,
		UserIDLegacy:  &userId,
		RatingID:      id,
		SourceTransID: orderNumber + "||product||Frtgffggffgft123||34343432",
		Value:         value,
		OrderNumber:   orderNumber,
		MediaPath: nil,
	}
	saveReq := []request.SaveRatingSubmissionMp{
		{
			UserID:        &userId,
			UserIDLegacy:  &userId,
			DisplayName:   &name,
			SourceTransID: orderNumber + "||product||Frtgffggffgft123||34343432",
			SourceUID:     "Frtgffggffgft123",
			SourceType:    "product",
			Value:         &value,
			MediaPath: nil,
			OrderNumber: orderNumber,
		},
	}
	arrSub := []entity.RatingSubmissionMp{sub}
	ratingMpRepository.Mock.On("FindRatingSubmissionBySourceTransID", SourceTransID).Return(nil, gorm.ErrRecordNotFound)
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
