package test

import (
	"context"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/message"
	"testing"

	publicresponse "go-klikdokter/app/model/response/public"

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
		Media:         []entity.MediaObj{},
		Comment:       userId,
		StoreUID:      "1",
	}

	objectID, _ := primitive.ObjectIDFromHex(id)
	ratingTypeID := entity.RatingTypesNumCol{
		ID: objectID,
	}
	sub := entity.RatingSubmissionMp{
		ID:            objectID,
		UserID:        &userId,
		UserIDLegacy:  &userId,
		RatingID:      id,
		SourceTransID: orderNumber + "||product||Frtgffggffgft123||34343432",
		Value:         value,
		OrderNumber:   orderNumber,
		Media:         nil,
		RatingTypeID:  ratingTypeID.ID.Hex(),
		Comment:       &userId,
	}

	saveReq := []entity.RatingSubmissionMp{
		{
			UserID:        &userId,
			UserIDLegacy:  &userId,
			DisplayName:   &name,
			SourceTransID: orderNumber + "||product||Frtgffggffgft123||34343432",
			SourceUID:     "Frtgffggffgft123",
			SourceType:    "product",
			Value:         value,
			Media:         nil,
			OrderNumber:   orderNumber,
			RatingTypeID:  ratingTypeID.ID.Hex(),
			Comment:       &userId,
			StoreUID:      input.StoreUID,
		},
	}
	arrSub := []entity.RatingSubmissionMp{sub}
	valueGroupBy := []publicresponse.PublicRatingSubGroupByValue{
		{
			ConvertedValue: 4,
			Total: 1,
		},
	}

	status := true

	formula := entity.RatingFormulaCol{
		ID: 	objectID,
		SourceType   :"product",
		Formula:    "(count/sum)/1" ,
		RatingTypeId : "123",
		RatingType: "rating_for_product",
		Status: &status, 
	}

	ratingMpRepository.Mock.On("FindRatingSubmissionBySourceTransID", SourceTransID).Return(nil, gorm.ErrRecordNotFound)
	ratingMpRepository.Mock.On("FindRatingTypeNumByRatingType", input.RatingType).Return(&ratingTypeID, nil)
	ratingMpRepository.Mock.On("CreateRatingSubmission", saveReq).Return(&arrSub, nil)
	ratingMpRepository.Mock.On("GetRatingSubsGroupByValue", input.SourceUID, "product").Return(valueGroupBy, nil)
	ratingMpRepository.Mock.On("GetRatingFormulaBySourceType", "product").Return(&formula, nil)

	_, msg := ratingMpSvc.CreateRatingSubmissionMp(context.Background(), input)

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
