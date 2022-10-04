package test

import (
	"encoding/json"
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util"
	"go-klikdokter/pkg/util/mocks"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var logger log.Logger

var ratingRepository = &repository_mock.RatingRepositoryMock{Mock: mock.Mock{}}
var medicalFacility = &mocks.MedicalFacilitySvc{Mock: mock.Mock{}}
var svc = service.NewRatingService(logger, ratingRepository, publicRatingRepository, medicalFacility)

func init() {
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}
}

var (
	value         = "4"
	id            = "629dce7bf1f26275e0d84826"
	description   = "alo"
	res           = "success"
	valueFailed   = "failed"
	Desc          = "Description"
	Bool          = true
	Scale         = 0
	ratingId      = "629ec07e6f3c2761ba2dc433"
	ratingTypeId  = "629ec07e6f3c2761ba2dc468"
	ratingtypeNum = "standard-0.0-to-5.0"
	name          = "name"
	statusTrue    = true
	sourceType    = "source type"
	callMFSuccess = "CallGetDetailMedicalFacilitySuccess"
	callMFFailed  = "CallGetDetailMedicalFacilityFailed"
	e             = errors.New("error")
)

var jwtObj = global.JWTObj{
	UserIdLegacy: "12345",
	Fullname:     "Test",
}

func TestCreateRatingTypeNum(t *testing.T) {
	var minScore = 0
	var scale = 0
	var maxScore = 5
	req := request.CreateRatingTypeNumRequest{
		Type:        "type",
		Description: &Desc,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Scale:       &scale,
		Intervals:   11,
		Status:      nil,
	}
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	ratingTypesNumCol := entity.RatingTypesNumCol{
		ID: objectId,
	}

	ratingRepository.Mock.On("CreateRatingTypeNum", req).Return(ratingTypesNumCol)

	_, msg := svc.CreateRatingTypeNum(req)
	assert.NotNil(t, message.SuccessMsg, msg)
}

func TestCreateRatingTypeNumErrScaleValueReq(t *testing.T) {
	var minScore = 0
	var scale = 3
	var maxScore = 5
	req := request.CreateRatingTypeNumRequest{
		Type:        "type",
		Description: &Desc,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Scale:       &scale,
		Intervals:   11,
		Status:      nil,
	}
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	ratingTypesNumCol := entity.RatingTypesNumCol{
		ID: objectId,
	}

	ratingRepository.Mock.On("CreateRatingTypeNum", req).Return(ratingTypesNumCol)

	_, msg := svc.CreateRatingTypeNum(req)
	assert.Equal(t, message.ErrScaleValueReq, msg)
}

func TestCreateRatingTypeNumErrSaveData(t *testing.T) {
	var minScore = 0
	var scale = 1
	var status bool
	var maxScore = 1
	req := request.CreateRatingTypeNumRequest{
		Type:        "12345",
		Description: &Desc,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Scale:       &scale,
		Intervals:   3,
		Status:      &status,
	}
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	ratingTypesNumCol := entity.RatingTypesNumCol{
		ID: objectId,
	}

	ratingRepository.Mock.On("CreateRatingTypeNum", req).Return(ratingTypesNumCol)

	_, msg := svc.CreateRatingTypeNum(req)
	assert.Equal(t, message.FailedMsg, msg)
}

func TestGetRatingTypeNumById(t *testing.T) {
	req := request.GetRatingTypeNumRequest{Id: "629ec0736f3c2761ba2dc867"}
	var minScore = 0
	var scale = 2
	var status bool
	var maxScore = 5
	var intervals = 11
	objectId, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")
	ratingTypesNumCol := entity.RatingTypesNumCol{
		ID:          objectId,
		Type:        "type",
		Description: &Desc,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Scale:       &scale,
		Intervals:   &intervals,
		Status:      &status,
	}

	ratingRepository.Mock.On("GetRatingTypeNumById", objectId).Return(ratingTypesNumCol)

	_, msg := svc.GetRatingTypeNumById(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetRatingTypeNumByIdFailed2(t *testing.T) {
	req := request.GetRatingTypeNumRequest{Id: "629dce7bf1f26275e0d84826"}
	var minScore = 0
	var scale = 2
	var status bool
	var maxScore = 2
	var intervals = 11

	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	ratingTypesNumCol := entity.RatingTypesNumCol{
		ID:          objectId,
		Type:        "type",
		Description: &Desc,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Scale:       &scale,
		Intervals:   &intervals,
		Status:      &status,
	}

	ratingRepository.Mock.On("GetRatingTypeNumById", objectId).Return(ratingTypesNumCol)

	_, msg := svc.GetRatingTypeNumById(req)
	assert.Equal(t, message.FailedMsg, msg)
}

func TestGetRatingTypeNumByIdErrNoData2(t *testing.T) {
	req := request.GetRatingTypeNumRequest{Id: "629ec07e6f3c2761ba2dc868"}

	objectId, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc868")

	ratingRepository.Mock.On("GetRatingTypeNumById", objectId).Return(nil)

	_, msg := svc.GetRatingTypeNumById(req)
	assert.Equal(t, message.ErrNoData, msg)
}

func TestGetRatingTypeNumByIdErrNoData1(t *testing.T) {
	req := request.GetRatingTypeNumRequest{Id: "213213213"}

	objectId, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc868")

	ratingRepository.Mock.On("GetRatingTypeNumById", objectId).Return(nil)

	_, msg := svc.GetRatingTypeNumById(req)
	assert.Equal(t, message.ErrNoData, msg)
}
func TestUpdateRatingTypeNum(t *testing.T) {
	var minScore = 0
	var scale = 0
	status := true
	var maxScore = 5
	var intervals = 6
	req := request.EditRatingTypeNumRequest{
		Id:          "629ec07e6f3c2761ba2dc867",
		Description: &Desc,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Scale:       &scale,
		Intervals:   &intervals,
		Status:      &status,
	}

	objectId, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc867")

	//rating := entity.RatingsCol{
	//	ID:             objectId,
	//	CommentAllowed: &status,
	//	Status:         &status,
	//}
	//objectId1, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")
	//rating := entity.RatingsCol{
	//	ID:   objectId1,
	//	Name: "abc",
	//}
	//objectId2, _ := primitive.ObjectIDFromHex("62a1bb7e0809e0a7bb12018b")
	//submissison := entity.RatingSubmisson{
	//	ID:      objectId2,
	//	Comment: "dlfslkjdf",
	//}
	objectId, _ = primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")

	ratingRepository.Mock.On("UpdateRatingTypeNum", objectId)
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(nil)
	//ratingRepository.Mock.On("GetRatingSubmissionByRatingId", rating.ID).Return(nil)

	msg := svc.UpdateRatingTypeNum(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestUpdateRatingTypeNumErrNodata(t *testing.T) {
	var minScore = 0
	var scale = 2
	var status bool
	var maxScore = 1
	var intervals = 4
	req := request.EditRatingTypeNumRequest{
		Id:          "23124",
		Description: &description,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Scale:       &scale,
		Intervals:   &intervals,
		Status:      &status,
	}
	objectId, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")

	ratingRepository.Mock.On("UpdateRatingTypeNum", objectId)
	msg := svc.UpdateRatingTypeNum(req)
	assert.Equal(t, message.ErrNoData, msg)
}

func TestUpdateRatingTypeNumFailed(t *testing.T) {
	var minScore = 0
	var scale = 0
	status := true
	var maxScore = 5
	var intervals = 4
	req := request.EditRatingTypeNumRequest{
		Id:          "629dce7bf1f26275e0d84826",
		Description: &Desc,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Scale:       &scale,
		Intervals:   &intervals,
		Status:      &status,
	}
	objectId1, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")
	rating := entity.RatingsCol{
		ID:   objectId1,
		Name: "abc",
	}
	//objectId2, _ := primitive.ObjectIDFromHex("62a1bb7e0809e0a7bb12018b")
	//submissison := entity.RatingSubmisson{
	//	ID:      objectId2,
	//	Comment: "dlfslkjdf",
	//}
	objectId, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")

	ratingRepository.Mock.On("UpdateRatingTypeNum", objectId)
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(rating)
	//ratingRepository.Mock.On("GetRatingSubmissionByRatingId", rating.ID).Return(nil)

	msg := svc.UpdateRatingTypeNum(req)
	assert.Equal(t, message.FailedMsg, msg)
}

func TestUpdateRatingTypeNumFailed2(t *testing.T) {
	var minScore = 0
	var scale = 2
	var status bool
	var maxScore = 1
	var intervals = 4
	req := request.EditRatingTypeNumRequest{
		Id:          "62a16b8afe7968dc56d6e47f",
		Description: &description,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Scale:       &scale,
		Intervals:   &intervals,
		Status:      &status,
	}
	objectId1, _ := primitive.ObjectIDFromHex("62a6e2ae8be76898e8ccc2e8")
	rating := entity.RatingsCol{
		ID:   objectId1,
		Name: "abc",
	}
	objectId2, _ := primitive.ObjectIDFromHex("62a6e2f78be76898e8ccc2e9")
	submissison := entity.RatingSubmisson{
		ID:       objectId2,
		Comment:  &Desc,
		RatingID: "62a6e2ae8be76898e8ccc2e8",
	}
	objectId, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")

	rateid := objectId1.Hex()
	ratingRepository.Mock.On("UpdateRatingTypeNum", objectId)
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(rating)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", rateid).Return(submissison)

	msg := svc.UpdateRatingTypeNum(req)
	assert.Equal(t, message.ErrCannotModifiedStatus, msg)
}

func TestUpdateRatingTypeNumInvalidIntervals(t *testing.T) {
	var minScore = 0
	var scale = 2
	var status bool
	var maxScore = 1
	var intervals = 22
	req := request.EditRatingTypeNumRequest{
		Id:          "629ec0736f3c2761ba2dc867",
		Description: &description,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Scale:       &scale,
		Intervals:   &intervals,
		Status:      &status,
	}
	objectId, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")

	ratingRepository.Mock.On("UpdateRatingTypeNum", objectId)
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(nil)
	//ratingRepository.Mock.On("GetRatingSubmissionByRatingId", rating.ID).Return(nil)

	msg := svc.UpdateRatingTypeNum(req)
	assert.Equal(t, message.ValidationFailCode, msg.Code)
}

func TestUpdateRatingTypeNumFailed3(t *testing.T) {
	var minScore = 0
	var scale = 2
	var status bool
	var maxScore = 1
	var intervals = 4
	req := request.EditRatingTypeNumRequest{
		Id:          "629ec07e6f3c2761ba2dc868",
		Description: &description,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Scale:       &scale,
		Intervals:   &intervals,
		Status:      &status,
	}
	//objectId1, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")
	//rating := entity.RatingsCol{
	//	ID:   objectId1,
	//	Name: "abc",
	//}
	//objectId2, _ := primitive.ObjectIDFromHex("62a1bb7e0809e0a7bb12018b")
	//submissison := entity.RatingSubmisson{
	//	ID:      objectId2,
	//	Comment: "dlfslkjdf",
	//}
	objectId, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")

	ratingRepository.Mock.On("UpdateRatingTypeNum", objectId)
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(nil)
	//ratingRepository.Mock.On("GetRatingSubmissionByRatingId", rating.ID).Return(nil)

	msg := svc.UpdateRatingTypeNum(req)
	assert.Equal(t, message.FailedMsg, msg)
}

//func TestUpdateRatingTypeNumErrIdFormatReq(t *testing.T) {
//	var minScore = 0
//	var scale = 2
//	var status bool
//	req := request.CreateRatingTypeNumRequest{
//		Id:          "sdkj234kld",
//		Type:        "12345",
//		Description: &Desc,
//		MinScore:    &minScore,
//		MaxScore:    1,
//		Scale:       &scale,
//		Intervals:   11,
//		Status:      &status,
//	}
//	objectId, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc868")
//
//	ratingRepository.Mock.On("UpdateRatingTypeNum", objectId)
//
//	msg := svc.UpdateRatingTypeNum(req)
//	assert.Equal(t, message.ErrNoData, msg)
//}

func TestDeleteRatingTypeNumById(t *testing.T) {
	req := request.GetRatingTypeNumRequest{Id: "629ec07e6f3c2761ba2dc867"}

	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(nil)
	msg := svc.DeleteRatingTypeNumById(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestDeleteRatingTypeNumByIdErrNoData(t *testing.T) {
	req := request.GetRatingTypeNumRequest{Id: "q324"}

	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(nil)
	msg := svc.DeleteRatingTypeNumById(req)
	assert.Equal(t, message.ErrNoData, msg)
}

func TestDeleteRatingTypeNumByIdFailed(t *testing.T) {
	req := request.GetRatingTypeNumRequest{Id: "629dce7bf1f26275e0d84826"}

	objectId1, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")
	rating := entity.RatingsCol{
		ID:   objectId1,
		Name: "abc",
	}
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(rating)
	msg := svc.DeleteRatingTypeNumById(req)
	assert.Equal(t, message.FailedMsg, msg)
}

func TestDeleteRatingTypeNumByIdErrThisRatingTypeIsInUse(t *testing.T) {
	req := request.GetRatingTypeNumRequest{Id: "629ec0736f3c2761ba2dc834"}
	objectId1, _ := primitive.ObjectIDFromHex("62a6e2ae8be76898e8ccc2e8")
	rating := entity.RatingsCol{
		ID:   objectId1,
		Name: "abc",
	}
	objectId2, _ := primitive.ObjectIDFromHex("62a6e2f78be76898e8ccc2e9")
	submissison := entity.RatingSubmisson{
		ID:       objectId2,
		Comment:  &Desc,
		RatingID: "62a6e2ae8be76898e8ccc2e8",
	}
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(rating)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", rating.ID).Return(submissison)
	msg := svc.DeleteRatingTypeNumById(req)
	assert.Equal(t, message.ErrThisRatingTypeIsInUse, msg)
}

func TestGetRatingTypeNums(t *testing.T) {
	req := request.GetRatingTypeNumsRequest{
		Sort:  "",
		Dir:   "desc",
		Page:  0,
		Limit: 0,
	}
	objectId1, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc468")
	objectId2, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc848")
	result := []entity.RatingTypesNumCol{
		{
			ID:          objectId1,
			Description: &description,
		},
		{
			ID:          objectId2,
			Description: &description,
		},
	}
	paginationResult := base.Pagination{
		Records:   120,
		Limit:     50,
		Page:      1,
		TotalPage: 12,
	}
	ratingRepository.Mock.On("GetRatingTypeNums", request.Filter{TypeId: []string(nil), MinScore: []int(nil), MaxScore: []int(nil)}, 1, int64(50), "updated_at", -1).Return(result, &paginationResult)

	_, _, msg := svc.GetRatingTypeNums(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetRatingSubmissionSuccess(t *testing.T) {
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	sub := entity.RatingSubmisson{
		ID:      objectId,
		Comment: &Desc,
	}

	ratingRepository.Mock.On("GetRatingSubmissionById", objectId).Return(sub, nil)

	_, msg := svc.GetRatingSubmission(id)

	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetRatingSubmissionFail(t *testing.T) {
	failId := "629dce7bf1f26275e0d84827"
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84827")

	ratingRepository.Mock.On("GetRatingSubmissionById", objectId).Return(nil, mongo.ErrNoDocuments)

	_, msg := svc.GetRatingSubmission(failId)

	assert.Equal(t, message.ErrRatingSubmissionNotFound, msg)
}

func TestGetRatingSubmissionWrongID(t *testing.T) {
	failId := "123456"
	objectId, _ := primitive.ObjectIDFromHex("123456")

	ratingRepository.Mock.On("GetRatingSubmissionById", objectId).Return(nil, mongo.ErrNoDocuments)

	_, msg := svc.GetRatingSubmission(failId)

	assert.Equal(t, message.ErrRatingSubmissionNotFound, msg)
}

func TestDeleteRatingSubmissionSuccess(t *testing.T) {
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")

	ratingRepository.Mock.On("DeleteSubmission", objectId).Return(nil)

	msg := svc.DeleteRatingSubmission(id)

	assert.Equal(t, message.SuccessMsg, msg)
}

func TestDeleteRatingSubmissionFail(t *testing.T) {
	failId := "629dce7bf1f26275e0d84827"
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84827")

	ratingRepository.Mock.On("DeleteSubmission", objectId).Return(errors.New("user not found"))

	msg := svc.DeleteRatingSubmission(failId)

	assert.Equal(t, message.ErrNoData, msg)
}

func TestDeleteRatingSubmissionWrongId(t *testing.T) {
	failId := "123456"
	objectId, _ := primitive.ObjectIDFromHex("123456")

	ratingRepository.Mock.On("DeleteSubmission", objectId).Return(nil)

	msg := svc.DeleteRatingSubmission(failId)

	assert.Equal(t, message.ErrDataNotFound, msg)
}

func TestCreateRatingSubmissionFailFindRating(t *testing.T) {
	valueRate := "1"
	likertID := "629dce7bf1f26275e0d84820"
	objectLikertId, _ := primitive.ObjectIDFromHex(likertID)
	input := request.CreateRatingSubmissionRequest{
		Ratings: []request.RatingByType{
			{
				ID:    likertID,
				Value: &valueRate,
			},
		},
		UserID:       &id,
		UserIDLegacy: &id,
		DisplayName:  &name,
	}

	sub := entity.RatingSubmisson{
		UserID:       &id,
		UserIDLegacy: &id,
		RatingID:     id,
	}

	likert := entity.RatingTypesLikertCol{
		ID:            objectLikertId,
		Description:   &description,
		NumStatements: 1,
		Statement01:   &description,
	}

	saveReq := []request.SaveRatingSubmission{
		{
			UserID:       &id,
			UserIDLegacy: &id,
			RatingID:     id,
		},
	}

	ratingRepository.Mock.On("FindRatingByRatingID", objectLikertId).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDLegacyAndRatingID", &id, id, id).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", &id, id, id).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingNumericTypeByRatingTypeID", objectLikertId).Return(nil, nil)
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", objectLikertId).Return(likert, nil)
	ratingRepository.Mock.On("CreateRatingSubmission", saveReq).Return(sub, nil)

	_, msg := svc.CreateRatingSubmission(input)

	assert.Equal(t, message.ErrRatingNotFound, msg)
}

func TestCreateRatingSubmissionFailAgentTooLong(t *testing.T) {
	valueRate := "1"
	likertID := "629dce7bf1f26275e0d84820"
	objectId, _ := primitive.ObjectIDFromHex(id)
	objectLikertId, _ := primitive.ObjectIDFromHex(likertID)
	input := request.CreateRatingSubmissionRequest{
		Ratings: []request.RatingByType{
			{
				ID:    id,
				Value: &valueRate,
			},
		},
		UserID:        &id,
		UserIDLegacy:  &id,
		SourceTransID: id,
		DisplayName:   &name,
		UserAgent:     "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
	}
	sourceTransIdConcenate := input.SourceTransID + "||" + id

	sub := entity.RatingSubmisson{
		UserID:        &id,
		UserIDLegacy:  &id,
		RatingID:      id,
		SourceTransID: "string",
	}

	rating := entity.RatingsCol{
		ID:             objectLikertId,
		RatingTypeId:   likertID,
		CommentAllowed: &Bool,
		Status:         &Bool,
	}

	likert := entity.RatingTypesLikertCol{
		ID:            objectLikertId,
		Description:   &description,
		NumStatements: 1,
		Statement01:   &description,
	}

	saveReq := []request.SaveRatingSubmission{
		{
			UserID:       &id,
			UserIDLegacy: &id,
			RatingID:     id,
		},
	}

	ratingRepository.Mock.On("FindRatingByRatingID", objectId).Return(rating, nil)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDLegacyAndRatingID", &id, id, sourceTransIdConcenate).Return(sub, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", &id, id, id).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingNumericTypeByRatingTypeID", objectLikertId).Return(nil, nil)
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", objectLikertId).Return(likert, nil)
	ratingRepository.Mock.On("CreateRatingSubmission", saveReq).Return(sub, nil)

	_, msg := svc.CreateRatingSubmission(input)

	assert.Equal(t, message.UserRated, msg)
}

func TestCreateRatingSubmissionFailLikertType(t *testing.T) {
	valueRate := "2"
	likertID := "629dce7bf1f26275e0d84820"
	objectId, _ := primitive.ObjectIDFromHex(id)
	objectLikertId, _ := primitive.ObjectIDFromHex(likertID)
	input := request.CreateRatingSubmissionRequest{
		Ratings: []request.RatingByType{
			{
				ID:    id,
				Value: &valueRate,
			},
		},
		UserID:       &id,
		UserIDLegacy: &id,
		DisplayName:  &name,
	}

	sub := entity.RatingSubmisson{
		UserID:       &id,
		UserIDLegacy: &id,
		RatingID:     id,
	}

	rating := entity.RatingsCol{
		ID:             objectId,
		RatingTypeId:   likertID,
		CommentAllowed: &Bool,
		Status:         &Bool,
	}

	likert := entity.RatingTypesLikertCol{
		ID:            objectLikertId,
		Description:   &description,
		NumStatements: 1,
		Statement01:   &description,
	}

	saveReq := []request.SaveRatingSubmission{
		{
			UserID:       &id,
			UserIDLegacy: &id,
			RatingID:     id,
		},
	}

	ratingRepository.Mock.On("FindRatingByRatingID", objectId).Return(rating, nil)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDLegacyAndRatingID", &id, id, id).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", &id, id, id).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingNumericTypeByRatingTypeID", objectLikertId).Return(nil, nil)
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", objectLikertId).Return(likert, nil)
	ratingRepository.Mock.On("CreateRatingSubmission", saveReq).Return(sub, nil)

	_, msg := svc.CreateRatingSubmission(input)

	assert.Equal(t, message.Message{
		Code:    message.ValidationFailCode,
		Message: "value must be integer and include in [1]",
	}, msg)
}

func TestCreateRatingSubmissionFailUidNull(t *testing.T) {
	valueRate := "1"
	likertID := "629dce7bf1f26275e0d84820"
	objectId, _ := primitive.ObjectIDFromHex(id)
	objectLikertId, _ := primitive.ObjectIDFromHex(likertID)
	input := request.CreateRatingSubmissionRequest{
		Ratings: []request.RatingByType{
			{
				ID:    id,
				Value: &valueRate,
			},
		},
		UserID:       nil,
		UserIDLegacy: nil,
	}

	sub := entity.RatingSubmisson{
		UserID:       &id,
		UserIDLegacy: &id,
		RatingID:     id,
	}

	rating := entity.RatingsCol{
		ID:             objectId,
		RatingTypeId:   likertID,
		CommentAllowed: &Bool,
		Status:         &Bool,
	}

	likert := entity.RatingTypesLikertCol{
		ID:            objectLikertId,
		Description:   &description,
		NumStatements: 1,
		Statement01:   &description,
	}

	saveReq := []request.SaveRatingSubmission{
		{
			UserID:       &id,
			UserIDLegacy: &id,
			RatingID:     id,
		},
	}

	ratingRepository.Mock.On("FindRatingByRatingID", objectId).Return(rating, nil)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDLegacyAndRatingID", &id, id, id).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", &id, id, id).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingNumericTypeByRatingTypeID", objectLikertId).Return(nil, nil)
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", objectLikertId).Return(likert, nil)
	ratingRepository.Mock.On("CreateRatingSubmission", saveReq).Return(sub, nil)

	_, msg := svc.CreateRatingSubmission(input)

	assert.Equal(t, message.UserUIDRequired, msg)
}

func TestUpdateRatingSubmissionSuccess(t *testing.T) {
	minScore := 0
	maxScore := 5
	intervals := 6
	valueRate := "1"
	numericID := "629dce7bf1f26275e0d84826"
	likertID := "629dce7bf1f26275e0d84826"
	objectId, _ := primitive.ObjectIDFromHex(id)
	objectIdT, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84827")
	input := request.UpdateRatingSubmissionRequest{
		ID:       id,
		Value:    &valueRate,
		RatingID: id,
	}

	sub := entity.RatingSubmisson{
		ID:            objectId,
		UserID:        &id,
		UserIDLegacy:  nil,
		RatingID:      id,
		SourceTransID: "id",
	}

	otherSub := entity.RatingSubmisson{
		ID:            objectIdT,
		UserID:        &id,
		UserIDLegacy:  nil,
		RatingID:      id,
		SourceTransID: "id1",
	}

	rating := entity.RatingsCol{
		ID:             objectId,
		RatingTypeId:   id,
		CommentAllowed: &Bool,
		Status:         &Bool,
	}

	num := entity.RatingTypesNumCol{
		ID:          objectId,
		Status:      &Bool,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Description: &Desc,
		Scale:       &Scale,
		Intervals:   &intervals,
	}

	objectNumericId, _ := primitive.ObjectIDFromHex(numericID)
	objectLikertId, _ := primitive.ObjectIDFromHex(likertID)
	ratingRepository.Mock.On("GetRatingSubmissionById", objectId).Return(sub, nil)
	ratingRepository.Mock.On("FindRatingByRatingID", objectId).Return(rating, nil)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", sub.UserID, input.RatingID, sub.SourceTransID).Return(otherSub, nil)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDLegacyAndRatingID", &res, "", "").Return()
	ratingRepository.Mock.On("FindRatingNumericTypeByRatingTypeID", objectNumericId).Return(num, nil)
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", objectLikertId).Return(num, nil)
	ratingRepository.Mock.On("UpdateRatingSubmission", input).Return(sub, nil)

	msg := svc.UpdateRatingSubmission(input)

	assert.Equal(t, message.SuccessMsg, msg)
}

func TestUpdateRatingSubmissionSuccessUserID(t *testing.T) {
	minScore := 0
	maxScore := 5
	intervals := 6
	valueRate := "1"
	numericID := "629dce7bf1f26275e0d84826"
	likertID := "629dce7bf1f26275e0d84826"
	objectId, _ := primitive.ObjectIDFromHex(id)
	input := request.UpdateRatingSubmissionRequest{
		ID:       id,
		Value:    &valueRate,
		RatingID: id,
	}

	sub := entity.RatingSubmisson{
		UserID:       &id,
		UserIDLegacy: nil,
		RatingID:     id,
	}

	rating := entity.RatingsCol{
		ID:             objectId,
		RatingTypeId:   id,
		CommentAllowed: &Bool,
		Status:         &Bool,
	}

	num := entity.RatingTypesNumCol{
		ID:          objectId,
		Status:      &Bool,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Description: &Desc,
		Scale:       &Scale,
		Intervals:   &intervals,
	}

	objectNumericId, _ := primitive.ObjectIDFromHex(numericID)
	objectLikertId, _ := primitive.ObjectIDFromHex(likertID)
	ratingRepository.Mock.On("GetRatingSubmissionById", objectId).Return(sub, nil)
	ratingRepository.Mock.On("FindRatingByRatingID", objectId).Return(rating, nil)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDLegacyAndRatingID", &id, id, "").Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", &id, id, id).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingNumericTypeByRatingTypeID", objectNumericId).Return(num, nil)
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", objectLikertId).Return(num, nil)
	ratingRepository.Mock.On("UpdateRatingSubmission", input).Return(sub, nil)

	msg := svc.UpdateRatingSubmission(input)

	assert.Equal(t, message.SuccessMsg, msg)
}

func TestUpdateRatingSubmissionSuccessUserIDLegacy(t *testing.T) {
	minScore := 0
	maxScore := 5
	intervals := 6
	valueRate := "1"
	numericID := "629dce7bf1f26275e0d84826"
	likertID := "629dce7bf1f26275e0d84826"
	objectId, _ := primitive.ObjectIDFromHex(id)
	input := request.UpdateRatingSubmissionRequest{
		ID:       id,
		Value:    &valueRate,
		RatingID: id,
	}

	sub := entity.RatingSubmisson{
		UserID:       &id,
		UserIDLegacy: &id,
		RatingID:     id,
	}

	rating := entity.RatingsCol{
		ID:             objectId,
		RatingTypeId:   id,
		CommentAllowed: &Bool,
		Status:         &Bool,
	}

	num := entity.RatingTypesNumCol{
		ID:          objectId,
		Status:      &Bool,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Description: &Desc,
		Scale:       &Scale,
		Intervals:   &intervals,
	}

	objectNumericId, _ := primitive.ObjectIDFromHex(numericID)
	objectLikertId, _ := primitive.ObjectIDFromHex(likertID)
	ratingRepository.Mock.On("GetRatingSubmissionById", objectId).Return(sub, nil)
	ratingRepository.Mock.On("FindRatingByRatingID", objectId).Return(rating, nil)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDLegacyAndRatingID", &id, id, id).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", &id, id, id).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingNumericTypeByRatingTypeID", objectNumericId).Return(num, nil)
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", objectLikertId).Return(num, nil)
	ratingRepository.Mock.On("UpdateRatingSubmission", input).Return(sub, nil)

	msg := svc.UpdateRatingSubmission(input)

	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetListRatingSubmission(t *testing.T) {
	matchStrValuePtr := "match"
	input := request.ListRatingSubmissionRequest{
		Dir:    "asc",
		Filter: "{\"user_uid\":[\"a12346fb-bd93-fedc-abcd-0739865540cb\",\"0739865540cb-bd93-fedc-abcd-a12346fb\"],\"score\":[4]}",
	}

	filter := request.RatingSubmissionFilter{}
	_ = json.Unmarshal([]byte(input.Filter), &filter)
	subs := []entity.RatingSubmisson{
		{
			UserID:       &matchStrValuePtr,
			UserIDLegacy: &matchStrValuePtr,
			Value:        value,
			Comment:      &Desc,
		},
	}

	page := base.Pagination{
		Records:      1,
		TotalRecords: 1,
		Limit:        50,
		Page:         1,
	}

	ratingRepository.Mock.On("GetListRatingSubmissions", filter, 1, int64(page.Limit), "created_at", 1).Return(subs, &page, nil)

	_, _, msg := svc.GetListRatingSubmissions(input)

	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetListRatingSubmissionMarshalErr(t *testing.T) {
	matchStrValuePtr := "match"
	input := request.ListRatingSubmissionRequest{
		Dir:    "asc",
		Filter: "{\"user_uid\":[\"a12346fb-bd93-fedc-abcd-0739865540cb\",\"0739865540cb-bd93-fedc-abcd-a12346fb\"],\"score\":[\"4\"]}",
	}

	filter := request.RatingSubmissionFilter{}
	_ = json.Unmarshal([]byte(input.Filter), &filter)
	subs := []entity.RatingSubmisson{
		{
			UserID:       &matchStrValuePtr,
			UserIDLegacy: &matchStrValuePtr,
			Value:        "4.5",
		},
	}

	page := base.Pagination{
		Records:      1,
		TotalRecords: 1,
		Limit:        50,
		Page:         1,
	}

	ratingRepository.Mock.On("GetListRatingSubmissions", filter, 1, int64(page.Limit), "created_at", 1).Return(subs, &page, nil)

	_, _, msg := svc.GetListRatingSubmissions(input)

	assert.Equal(t, message.WrongFilter, msg)
}

func TestCreateRatingTypeLikert(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Type:        "type",
		Description: &description,
	}
	ratingRepository.Mock.On("CreateRatingTypeLikert", req).Return()

	msg := svc.CreateRatingTypeLikert(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestCreateRatingTypeLikertFailed(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Type:        "typeErr",
		Description: &description,
	}
	ratingRepository.Mock.On("CreateRatingTypeLikert", req).Return()

	msg := svc.CreateRatingTypeLikert(req)
	assert.Equal(t, message.FailedMsg, msg)
}

func TestCreateRatingTypeLikertFailed2(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Type:        "duplicate",
		Description: &description,
	}
	ratingRepository.Mock.On("CreateRatingTypeLikert", req).Return()

	msg := svc.CreateRatingTypeLikert(req)
	assert.Equal(t, message.ErrDuplicateType, msg)
}

func TestCreateRatingTypeLikertErrSaveData(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:            "629ec07e6f3c2761ba2dc828",
		NumStatements: 1,
		Statement01:   &description,
		Statement02:   &description,
	}

	msg := svc.CreateRatingTypeLikert(req)
	assert.Equal(t, message.ErrMatchNumState, msg)
}

func TestCreateRatingTypeLikertErrSaveData2(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:            "629ec07e6f3c2761ba2dc828",
		NumStatements: 3,
		Statement01:   &description,
		Statement02:   &description,
	}
	msg := svc.CreateRatingTypeLikert(req)
	assert.Equal(t, message.ErrMatchNumState, msg)
}

func TestGetRatingTypeLikertById(t *testing.T) {
	req := request.GetRatingTypeLikertRequest{
		Id: "629ec07e6f3c2761ba2dc868",
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)

	likert := entity.RatingTypesLikertCol{
		Type:        "test",
		Description: &description,
	}
	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(likert)

	_, msg := svc.GetRatingTypeLikertById(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetRatingTypeLikertByIdErrIdFormatReq(t *testing.T) {
	req := request.GetRatingTypeLikertRequest{
		Id: "3411ds",
	}
	_, msg := svc.GetRatingTypeLikertById(req)
	assert.Equal(t, message.ErrNoData, msg)
}

func TestGetRatingTypeLikertByIdFailed(t *testing.T) {
	req := request.GetRatingTypeLikertRequest{
		Id: "629dce7bf1f26275e0d84826",
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)

	likert := entity.RatingTypesLikertCol{
		Type:        "test",
		Description: &description,
	}
	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(likert)

	_, msg := svc.GetRatingTypeLikertById(req)
	assert.Equal(t, message.FailedMsg, msg)
}

func TestGetRatingTypeLikertByIdFailedErrNoData(t *testing.T) {
	req := request.GetRatingTypeLikertRequest{
		Id: "629dce7bf1f26275e0d84326",
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)
	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(nil)

	_, msg := svc.GetRatingTypeLikertById(req)
	assert.Equal(t, message.ErrNoData, msg)
}

func TestUpdateRatingTypeLikertSuccess(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:          "629ec07e6f3c2761ba2dc828",
		Description: &description,
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)
	likert := entity.RatingTypesLikertCol{
		Description: &description,
	}
	rating := entity.RatingsCol{
		ID:           objectId,
		RatingTypeId: "629dce7bf1f26275e0d84824",
	}
	submission := entity.RatingSubmisson{
		ID: objectId,
	}
	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(likert)
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(rating)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", req.Id).Return(submission)
	ratingRepository.Mock.On("UpdateRatingTypeLikert", objectId, req).Return()

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestUpdateRatingTypeLikertSuccess2(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:          "629ec07e6f3c2761ba2dc828",
		Description: &description,
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)
	likert := entity.RatingTypesLikertCol{
		Description: &description,
	}
	rating := entity.RatingsCol{
		ID:           objectId,
		RatingTypeId: "629dce7bf1f26275e0d84824",
	}
	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(likert)
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(rating)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", req.Id).Return(nil)
	ratingRepository.Mock.On("UpdateRatingTypeLikert", objectId, req).Return()

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestUpdateRatingTypeLikertSuccess3(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:          "629dce7bf1f26275e0d84821",
		Description: &description,
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)
	likert := entity.RatingTypesLikertCol{
		Description: &description,
	}

	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(likert)
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(nil)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId").Return(nil)
	ratingRepository.Mock.On("UpdateRatingTypeLikert", objectId, req).Return()

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestUpdateRatingTypeLikertErrSaveData(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:            "629ec07e6f3c2761ba2dc821",
		NumStatements: 1,
		Statement01:   &description,
		Statement02:   &description,
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)

	likert := entity.RatingTypesLikertCol{
		Type:          "test",
		NumStatements: 2,
	}
	rating := entity.RatingsCol{
		ID: objectId,
	}

	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(likert)
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(rating)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", rating.ID.Hex()).Return(nil)

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.ErrMatchNumState, msg)
}

func TestUpdateRatingTypeLikertErrSaveData2(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:            "629ec07e6f3c2761ba2dc821",
		NumStatements: 3,
		Statement01:   &description,
		Statement02:   &description,
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)

	likert := entity.RatingTypesLikertCol{
		Type:          "test",
		NumStatements: 2,
	}
	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(likert)
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(nil)

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.ErrMatchNumState, msg)
}

func TestUpdateRatingTypeLikertFailed3(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:          "629dce7bf1f26275e0d84826",
		Description: &valueFailed,
	}
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	likert := entity.RatingTypesLikertCol{
		Type:          "test",
		NumStatements: 2,
	}
	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(likert)
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(nil)

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.FailedMsg, msg)
}

func TestUpdateRatingTypeLikertFailedRequired1(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:          "629ec07e6f3c2761ba2dc824",
		Type:        "123",
		Description: &description,
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)
	likert := entity.RatingTypesLikertCol{
		Type:        "123",
		Description: &description,
	}
	rating := entity.RatingsCol{
		ID:           objectId,
		RatingTypeId: "62a6e2ae8be76898e8ccc2e8",
	}
	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(likert)
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(rating)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", rating.ID.Hex()).Return(nil)
	ratingRepository.Mock.On("UpdateRatingTypeLikert", objectId, req).Return()

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.ErrCannotModifiedType, msg)
}

func TestUpdateRatingTypeLikertWrongTypeID(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:          "A",
		Type:        "123",
		Description: &description,
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)
	likert := entity.RatingTypesLikertCol{
		Type:        "A",
		Description: &description,
	}
	rating := entity.RatingsCol{
		ID:           objectId,
		RatingTypeId: "629dce7bf1f26275e0d84824",
	}
	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(likert)
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(rating)
	ratingRepository.Mock.On("UpdateRatingTypeLikert", objectId, req).Return()

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.ErrNoData, msg)
}

func TestDeleteRatingTypeLikertById(t *testing.T) {
	req := request.GetRatingTypeLikertRequest{
		Id: "629ec07e6f3c2761ba2dc868",
	}
	msg := svc.DeleteRatingTypeLikertById(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestDeleteRatingTypeLikertByIdErrIdFormatReq(t *testing.T) {
	req := request.GetRatingTypeLikertRequest{
		Id: "21323",
	}
	msg := svc.DeleteRatingTypeLikertById(req)
	assert.Equal(t, message.ErrNoData, msg)
}

func TestDeleteRatingTypeLikertByIdFailedMsg(t *testing.T) {
	req := request.GetRatingTypeLikertRequest{Id: "629dce7bf1f26275e0d84826"}

	objectId1, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")
	rating := entity.RatingsCol{
		ID:   objectId1,
		Name: "abc",
	}
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(rating)
	msg := svc.DeleteRatingTypeLikertById(req)
	assert.Equal(t, message.FailedMsg, msg)
}

func TestDeleteRatingTypeLikertByIdErrNoData(t *testing.T) {
	req := request.GetRatingTypeLikertRequest{Id: "q324"}

	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(nil)
	msg := svc.DeleteRatingTypeLikertById(req)
	assert.Equal(t, message.ErrNoData, msg)
}

func TestDeleteRatingTypeLikertByIdErrThisRatingTypeIsInUse(t *testing.T) {
	req := request.GetRatingTypeLikertRequest{Id: "629ec0736f3c2761ba2dc834"}
	objectId1, _ := primitive.ObjectIDFromHex("62a6e2ae8be76898e8ccc2e8")
	rating := entity.RatingsCol{
		ID:   objectId1,
		Name: "abc",
	}
	objectId2, _ := primitive.ObjectIDFromHex("62a6e2f78be76898e8ccc2e9")
	submissison := entity.RatingSubmisson{
		ID:       objectId2,
		Comment:  &Desc,
		RatingID: "62a6e2ae8be76898e8ccc2e8",
	}
	ratingRepository.Mock.On("GetRatingByType", req.Id).Return(rating)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", rating.ID).Return(submissison)
	msg := svc.DeleteRatingTypeLikertById(req)
	assert.Equal(t, message.ErrThisRatingTypeIsInUse, msg)
}

func TestGetRatingTypeLikerts(t *testing.T) {
	req := request.GetRatingTypeLikertsRequest{
		Sort:  "",
		Dir:   "desc",
		Page:  0,
		Limit: 0,
	}
	objectId1, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc468")
	objectId2, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc848")
	result := []entity.RatingTypesLikertCol{
		{
			ID:          objectId1,
			Description: &description,
		},
		{
			ID:          objectId2,
			Description: &description,
		},
	}
	paginationResult := base.Pagination{
		Records:   120,
		Limit:     50,
		Page:      1,
		TotalPage: 12,
	}
	ratingRepository.Mock.On("GetRatingTypeLikerts", request.FilterRatingTypeLikert{TypeId: []string(nil)}, 1, int64(50), "updated_at", -1).Return(result, &paginationResult)

	_, _, msg := svc.GetRatingTypeLikerts(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetRatingTypeLikertsFailedMsg(t *testing.T) {
	req := request.GetRatingTypeLikertsRequest{
		Sort:  "failed",
		Dir:   "asc",
		Page:  0,
		Limit: 0,
	}
	objectId1, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc468")
	objectId2, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc848")
	result := []entity.RatingTypesLikertCol{
		{
			ID:          objectId1,
			Description: &description,
		},
		{
			ID:          objectId2,
			Description: &description,
		},
	}
	paginationResult := base.Pagination{
		Records:   120,
		Limit:     50,
		Page:      1,
		TotalPage: 12,
	}
	ratingRepository.Mock.On("GetRatingTypeLikerts", request.FilterRatingTypeLikert{TypeId: []string(nil)}, 1, int64(50), "failed", 1).Return(result, &paginationResult)

	_, _, msg := svc.GetRatingTypeLikerts(req)
	assert.Equal(t, message.FailedMsg, msg)
}

func TestCreateRating(t *testing.T) {
	req := request.SaveRatingRequest{
		Name:         name + "1",
		SourceUid:    callMFSuccess,
		SourceType:   sourceType + "1",
		RatingType:   ratingtypeNum,
		RatingTypeId: ratingTypeId,
		Status:       &statusTrue,
	}

	rating := &entity.RatingsCol{
		Name:   name + "1",
		Status: &statusTrue,
	}

	ratingTypeNum := &entity.RatingTypesNumCol{
		Type: ratingtypeNum,
	}

	ObjRatingTypeId, _ := primitive.ObjectIDFromHex(req.RatingTypeId)

	medicalFacilityResponseHttp := &util.ResponseHttp{
		Meta: util.MetaResponse{
			Code:    200,
			Message: "OK",
		},
		Data: util.Data{},
	}

	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.SourceUid, req.SourceType).Return(nil, nil)
	medicalFacility.Mock.On("CallGetDetailMedicalFacility", req.SourceUid).Return(medicalFacilityResponseHttp, nil)
	ratingRepository.Mock.On("GetRatingTypeNumByIdAndStatus", ObjRatingTypeId).Return(ratingTypeNum, nil)
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", ObjRatingTypeId).Return(nil, nil)
	ratingRepository.Mock.On("CreateRating", req).Return(rating, nil)

	_, mgs := svc.CreateRating(req)
	assert.Equal(t, message.SuccessMsg, mgs)
}

func TestCreateRatingErrRatingTypeNotExist2(t *testing.T) {
	req := request.SaveRatingRequest{
		Name:         name + "1",
		SourceUid:    callMFSuccess,
		SourceType:   sourceType + "1",
		RatingType:   ratingtypeNum,
		RatingTypeId: "62a950cb46b3c9f96df11bde",
		Status:       nil,
	}

	rating := &entity.RatingsCol{
		Name:   name + "1",
		Status: &statusTrue,
	}

	ratingTypeNum := &entity.RatingTypesNumCol{
		Type: ratingtypeNum,
	}

	ObjRatingTypeId, _ := primitive.ObjectIDFromHex(req.RatingTypeId)

	medicalFacilityResponseHttp := &util.ResponseHttp{
		Meta: util.MetaResponse{
			Code:    200,
			Message: "OK",
		},
		Data: util.Data{},
	}

	ratingTypeLikert := &entity.RatingTypesLikertCol{
		Type: "testfailed1",
	}
	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.SourceUid, req.SourceType).Return(nil, nil)
	medicalFacility.Mock.On("CallGetDetailMedicalFacility", req.SourceUid).Return(medicalFacilityResponseHttp, nil)
	ratingRepository.Mock.On("GetRatingTypeNumByIdAndStatus", ObjRatingTypeId).Return(ratingTypeNum, nil)
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", ObjRatingTypeId).Return(ratingTypeLikert, nil)
	ratingRepository.Mock.On("CreateRating", req).Return(rating, nil)

	_, mgs := svc.CreateRating(req)
	assert.Equal(t, message.ErrRatingTypeNotExist, mgs)
}

func TestCreateRatingErrRatingTypeNotExist1(t *testing.T) {
	req := request.SaveRatingRequest{
		Name:         name + "1",
		SourceUid:    callMFSuccess,
		SourceType:   sourceType + "1",
		RatingType:   ratingtypeNum,
		RatingTypeId: "62a950cb46b3c9f96df11bda",
		Status:       &statusTrue,
	}

	rating := &entity.RatingsCol{
		Name:   name + "1",
		Status: &statusTrue,
	}

	ObjRatingTypeId, _ := primitive.ObjectIDFromHex(req.RatingTypeId)

	medicalFacilityResponseHttp := &util.ResponseHttp{
		Meta: util.MetaResponse{
			Code:    200,
			Message: "OK",
		},
		Data: util.Data{},
	}

	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.SourceUid, req.SourceType).Return(nil, nil)
	medicalFacility.Mock.On("CallGetDetailMedicalFacility", req.SourceUid).Return(medicalFacilityResponseHttp, nil)
	ratingRepository.Mock.On("GetRatingTypeNumByIdAndStatus", ObjRatingTypeId).Return(nil, nil)
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", ObjRatingTypeId).Return(nil, nil)
	ratingRepository.Mock.On("CreateRating", req).Return(rating, nil)

	_, mgs := svc.CreateRating(req)
	assert.Equal(t, message.ErrRatingTypeNotExist, mgs)
}

func TestCreateRatingFailed1(t *testing.T) {
	req := request.SaveRatingRequest{
		Name:         name + "2",
		Description:  nil,
		SourceUid:    callMFFailed,
		SourceType:   "hospital",
		RatingType:   ratingtypeNum,
		RatingTypeId: ratingTypeId,
	}

	medicalFacilityResponseHttp := &util.ResponseHttp{
		Meta: util.MetaResponse{
			Code:    400,
			Message: "Data tidak ditemukan",
		},
		Data: util.Data{},
	}

	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.SourceUid, req.SourceType).Return(nil, nil)
	medicalFacility.Mock.On("CallGetDetailMedicalFacility", req.SourceUid).Return(medicalFacilityResponseHttp, nil)

	_, mgs := svc.CreateRating(req)
	assert.Equal(t, message.ErrSourceNotExist, mgs)
}

func TestCreateRatingFailed3(t *testing.T) {
	req := request.SaveRatingRequest{
		Name:         name + "3",
		SourceUid:    callMFSuccess,
		SourceType:   "hospital",
		RatingType:   ratingtypeNum,
		RatingTypeId: "testFailed",
	}

	medicalFacilityResponseHttp := &util.ResponseHttp{
		Meta: util.MetaResponse{
			Code:    200,
			Message: "OK",
		},
		Data: util.Data{},
	}

	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.SourceUid, req.SourceType).Return(nil, nil)
	medicalFacility.Mock.On("CallGetDetailMedicalFacility", req.SourceUid).Return(medicalFacilityResponseHttp, nil)
	_, mgs := svc.CreateRating(req)
	assert.Equal(t, message.ErrRatingTypeNotExist, mgs)
}

func TestCreateRatingFailed5(t *testing.T) {
	req := request.SaveRatingRequest{
		Name:         name + "5",
		SourceUid:    callMFSuccess,
		SourceType:   "hospital",
		RatingTypeId: "629ec07e6f3c2761ba2dc461",
	}

	ObjRatingTypeId, _ := primitive.ObjectIDFromHex(req.RatingTypeId)

	medicalFacilityResponseHttp := &util.ResponseHttp{
		Meta: util.MetaResponse{
			Code:    200,
			Message: "OK",
		},
		Data: util.Data{},
	}

	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.SourceUid, req.SourceType).Return(nil, nil)
	medicalFacility.Mock.On("CallGetDetailMedicalFacility", req.SourceUid).Return(medicalFacilityResponseHttp, nil)
	ratingRepository.Mock.On("GetRatingTypeNumByIdAndStatus", ObjRatingTypeId).Return(nil, e)

	_, mgs := svc.CreateRating(req)
	assert.Equal(t, message.FailedMsg, mgs)
}

func TestCreateRatingFailed6(t *testing.T) {
	req := request.SaveRatingRequest{
		Name:         name + "6",
		SourceUid:    callMFSuccess,
		SourceType:   "hospital",
		RatingTypeId: "629ec07e6f3c2761ba2dc462",
	}

	ratingTypeNum := &entity.RatingTypesNumCol{
		Type: ratingtypeNum,
	}

	ObjRatingTypeId, _ := primitive.ObjectIDFromHex(req.RatingTypeId)

	medicalFacilityResponseHttp := &util.ResponseHttp{
		Meta: util.MetaResponse{
			Code:    200,
			Message: "OK",
		},
		Data: util.Data{},
	}

	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.SourceUid, req.SourceType).Return(nil, nil)
	medicalFacility.Mock.On("CallGetDetailMedicalFacility", req.SourceUid).Return(medicalFacilityResponseHttp, nil)
	ratingRepository.Mock.On("GetRatingTypeNumByIdAndStatus", ObjRatingTypeId).Return(ratingTypeNum, nil)
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", ObjRatingTypeId).Return(nil, e)

	_, mgs := svc.CreateRating(req)
	assert.Equal(t, message.FailedMsg, mgs)
}

func TestCreateRatingFailed7(t *testing.T) {
	req := request.SaveRatingRequest{
		Name:         name + "7",
		SourceUid:    callMFSuccess,
		SourceType:   "hospital",
		RatingTypeId: ratingTypeId,
	}

	ObjRatingTypeId, _ := primitive.ObjectIDFromHex(req.RatingTypeId)

	medicalFacilityResponseHttp := &util.ResponseHttp{
		Meta: util.MetaResponse{
			Code:    200,
			Message: "OK",
		},
		Data: util.Data{},
	}

	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.SourceUid, req.SourceType).Return(nil, nil)
	medicalFacility.Mock.On("CallGetDetailMedicalFacility", req.SourceUid).Return(medicalFacilityResponseHttp, nil)
	ratingRepository.Mock.On("GetRatingTypeNumByIdAndStatus", ObjRatingTypeId).Return(nil, nil)
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", ObjRatingTypeId).Return(nil, nil)

	_, mgs := svc.CreateRating(req)
	assert.Equal(t, message.ErrRatingTypeNotExist, mgs)
}

func TestCreateRatingFailed8(t *testing.T) {
	req := request.SaveRatingRequest{
		Name: name + "68",
	}

	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.SourceUid, req.SourceType).Return(nil, e)

	_, mgs := svc.CreateRating(req)
	assert.Equal(t, message.FailedMsg, mgs)
}

func TestCreateRatingFailed9(t *testing.T) {
	req := request.SaveRatingRequest{
		Name:       name + "69",
		SourceUid:  "69",
		SourceType: "69",
	}

	rating := &entity.RatingsCol{
		SourceUid:  "69",
		SourceType: "69",
	}

	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.SourceUid, req.SourceType).Return(rating, nil)

	_, mgs := svc.CreateRating(req)
	assert.Equal(t, message.ErrExistingRatingTypeIdSourceUidAndSourceType, mgs)
}

func TestUpdateRating(t *testing.T) {
	req := request.UpdateRatingRequest{
		Id: ratingId,
		Body: request.BodyUpdateRatingRequest{
			Name:       name + "11",
			SourceUid:  callMFSuccess,
			SourceType: sourceType + "11",
		},
	}

	ObjRatingId, _ := primitive.ObjectIDFromHex(req.Id)

	rating := &entity.RatingsCol{
		Name: name + "11",
	}

	medicalFacilityResponseHttp := &util.ResponseHttp{
		Meta: util.MetaResponse{
			Code:    200,
			Message: "OK",
		},
		Data: util.Data{},
	}

	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.Body.SourceUid, req.Body.SourceType).Return(nil, nil)
	medicalFacility.Mock.On("CallGetDetailMedicalFacility", req.Body.SourceUid).Return(medicalFacilityResponseHttp, nil)
	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(rating, nil)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", req.Id).Once().Return(nil, nil)
	ratingRepository.Mock.On("UpdateRating", ObjRatingId, req.Body).Return(rating, nil)

	mgs := svc.UpdateRating(req)
	assert.Equal(t, message.SuccessMsg, mgs)
}

func TestUpdateRatingFailed(t *testing.T) {
	req := request.UpdateRatingRequest{
		Id: ratingId,
		Body: request.BodyUpdateRatingRequest{
			Name:       name + "21",
			SourceUid:  callMFSuccess,
			SourceType: sourceType + "21",
		},
	}

	ObjRatingId, _ := primitive.ObjectIDFromHex(req.Id)

	rating := &entity.RatingsCol{
		Name: name + "11",
	}

	medicalFacilityResponseHttp := &util.ResponseHttp{
		Meta: util.MetaResponse{
			Code:    200,
			Message: "OK",
		},
		Data: util.Data{},
	}

	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.Body.SourceUid, req.Body.SourceType).Return(nil, nil)
	medicalFacility.Mock.On("CallGetDetailMedicalFacility", req.Body.SourceUid).Return(medicalFacilityResponseHttp, nil)
	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(rating, nil)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", req.Id).Return(nil, nil)
	ratingRepository.Mock.On("UpdateRating", ObjRatingId, req.Body).Return(nil, e)

	mgs := svc.UpdateRating(req)
	assert.Equal(t, message.FailedMsg, mgs)
}

func TestUpdateRatingFailed1(t *testing.T) {
	req := request.UpdateRatingRequest{
		Id: ratingId,
		Body: request.BodyUpdateRatingRequest{
			Name:       name + "12",
			SourceUid:  callMFFailed,
			SourceType: sourceType + "12",
		},
	}

	rating := &entity.RatingsCol{
		Name: name + "12",
	}

	ObjRatingId, _ := primitive.ObjectIDFromHex(req.Id)

	medicalFacilityResponseHttp := &util.ResponseHttp{
		Meta: util.MetaResponse{
			Code:    400,
			Message: "Data tidak ditemukan",
		},
		Data: util.Data{},
	}

	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.Body.SourceUid, req.Body.SourceType).Return(nil, nil)
	medicalFacility.Mock.On("CallGetDetailMedicalFacility", req.Body.SourceUid).Return(medicalFacilityResponseHttp, nil)
	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(rating, nil)
	mgs := svc.UpdateRating(req)
	assert.Equal(t, message.ErrSourceNotExist, mgs)
}

func TestUpdateRatingFailed2(t *testing.T) {
	req := request.UpdateRatingRequest{
		Id: ratingId,
		Body: request.BodyUpdateRatingRequest{
			Name:       name + "43",
			SourceUid:  callMFSuccess,
			SourceType: sourceType + "43",
		},
	}

	ObjRatingId, _ := primitive.ObjectIDFromHex(req.Id)

	rating := &entity.RatingsCol{
		Name: name + "13",
	}

	medicalFacilityResponseHttpSuccess := &util.ResponseHttp{
		Meta: util.MetaResponse{
			Code:    200,
			Message: "OK",
		},
		Data: util.Data{},
	}

	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(rating, nil)
	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.Body.SourceUid, req.Body.SourceType).Return(nil, e)
	medicalFacility.Mock.On("CallGetDetailMedicalFacility", req.Body.SourceUid).Return(medicalFacilityResponseHttpSuccess, nil)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", req.Id).Return(nil, nil)
	mgs := svc.UpdateRating(req)
	assert.Equal(t, message.FailedMsg, mgs)
}

func TestUpdateRatingFailed3(t *testing.T) {
	req := request.UpdateRatingRequest{
		Id: ratingId,
		Body: request.BodyUpdateRatingRequest{
			Name:       name + "14",
			SourceUid:  callMFSuccess,
			SourceType: sourceType + "14",
		},
	}

	ObjRatingId, _ := primitive.ObjectIDFromHex(req.Id)

	rating := &entity.RatingsCol{
		Name: name + "14",
	}

	medicalFacilityResponseHttp := &util.ResponseHttp{
		Meta: util.MetaResponse{
			Code:    200,
			Message: "OK",
		},
		Data: util.Data{},
	}

	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(rating, nil)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", req.Id).Return(nil, nil)
	medicalFacility.Mock.On("CallGetDetailMedicalFacility", req.Body.SourceUid).Return(medicalFacilityResponseHttp, nil)
	ratingRepository.Mock.On("GetRatingByRatingTypeSourceUidAndSourceType", req.Body.SourceUid, req.Body.SourceType).Return(rating, nil)
	mgs := svc.UpdateRating(req)
	assert.Equal(t, message.ErrExistingRatingTypeIdSourceUidAndSourceType, mgs)
}

func TestUpdateRatingFailed4(t *testing.T) {
	req := request.UpdateRatingRequest{
		Id: "629dce7bf1f26275e0d84826",
		Body: request.BodyUpdateRatingRequest{
			Name: name + "55",
		},
	}

	ObjRatingId, _ := primitive.ObjectIDFromHex(req.Id)

	rating := &entity.RatingsCol{
		ID:   ObjRatingId,
		Name: name + "23",
	}

	ratingSubmission := entity.RatingSubmisson{
		RatingID: ratingId,
	}

	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(rating, nil)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", "629dce7bf1f26275e0d84826").Return(ratingSubmission)
	mgs := svc.UpdateRating(req)
	assert.Equal(t, message.FailedMsg, mgs)
}

func TestUpdateRatingFailed8(t *testing.T) {
	req := request.UpdateRatingRequest{
		Id: "testUpdateFailed",
	}

	mgs := svc.UpdateRating(req)
	assert.Equal(t, message.ErrDataNotFound, mgs)
}

func TestUpdateRatingFailed9(t *testing.T) {
	req := request.UpdateRatingRequest{
		Id: "629ec07e6f3c2761ba2dc411",
	}

	ObjRatingId, _ := primitive.ObjectIDFromHex(req.Id)

	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(nil, e)

	mgs := svc.UpdateRating(req)
	assert.Equal(t, message.FailedMsg, mgs)
}

func TestUpdateRatingFailed10(t *testing.T) {
	req := request.UpdateRatingRequest{
		Id: "629ec07e6f3c2761ba2dc422",
	}

	ObjRatingId, _ := primitive.ObjectIDFromHex(req.Id)

	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(nil, mongo.ErrNoDocuments)

	mgs := svc.UpdateRating(req)
	assert.Equal(t, message.ErrDataNotFound, mgs)
}

func TestGetRatingById(t *testing.T) {
	rating := &entity.RatingsCol{
		Name: name + "21",
	}

	ObjRatingId, _ := primitive.ObjectIDFromHex(ratingId)

	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(rating, nil)

	result, mgs := svc.GetRatingById(ratingId)
	assert.Equal(t, message.SuccessMsg, mgs)
	assert.NotNil(t, result)
}

func TestGetRatingByIdFailed1(t *testing.T) {
	ratingIdFailed := "629ec07e6f3c2761ba2dc111"

	ObjRatingId, _ := primitive.ObjectIDFromHex(ratingIdFailed)

	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(nil, e)
	result, mgs := svc.GetRatingById(ratingIdFailed)
	assert.Equal(t, message.FailedMsg, mgs)
	assert.Nil(t, result)
}

func TestGetRatingByIdFailed2(t *testing.T) {
	ratingIdFailed := "getRatingFailed"

	result, mgs := svc.GetRatingById(ratingIdFailed)
	assert.Equal(t, message.ErrDataNotFound, mgs)
	assert.Nil(t, result)
}

func TestGetRatingByIdFailed3(t *testing.T) {
	ratingIdFailed := "629ec07e6f3c2761ba2dc112"

	ObjRatingId, _ := primitive.ObjectIDFromHex(ratingIdFailed)

	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(nil, mongo.ErrNoDocuments)
	result, mgs := svc.GetRatingById(ratingIdFailed)
	assert.Equal(t, message.ErrDataNotFound, mgs)
	assert.Nil(t, result)
}

func TestDeleteRating(t *testing.T) {
	rId := "629ec07e6f3c2761ba2dc414"
	ObjRatingId, _ := primitive.ObjectIDFromHex(rId)

	rating := &entity.RatingsCol{
		ID:   ObjRatingId,
		Name: name + "22",
	}

	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(rating, nil)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", rId).Return(nil)
	ratingRepository.Mock.On("DeleteRating", ObjRatingId).Return(nil)

	mgs := svc.DeleteRating(rId)
	assert.Equal(t, message.SuccessMsg, mgs)
}

func TestDeleteRatingFailed8(t *testing.T) {
	rId := "629ec07e6f3c2761ba2dc412"
	ObjRatingId, _ := primitive.ObjectIDFromHex(rId)

	rating := &entity.RatingsCol{
		ID:   ObjRatingId,
		Name: name + "22",
	}

	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(rating, nil)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", rId).Return(nil)
	ratingRepository.Mock.On("DeleteRating", ObjRatingId).Return(e)

	mgs := svc.DeleteRating(rId)
	assert.Equal(t, message.FailedMsg, mgs)
}

func TestDeleteRatingFailed1(t *testing.T) {
	ratingIdFailed := "629ec07e6f3c2761ba2dc111"

	ObjRatingId, _ := primitive.ObjectIDFromHex(ratingIdFailed)

	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(nil, e)
	mgs := svc.DeleteRating(ratingIdFailed)
	assert.Equal(t, message.FailedMsg, mgs)
}

func TestDeleteRatingFailed2(t *testing.T) {
	ratingIdFailed := "getRatingFailed"

	mgs := svc.DeleteRating(ratingIdFailed)
	assert.Equal(t, message.ErrDataNotFound, mgs)
}

func TestDeleteRatingFailed3(t *testing.T) {
	ratingIdFailed := "629ec07e6f3c2761ba2dc112"

	ObjRatingId, _ := primitive.ObjectIDFromHex(ratingIdFailed)

	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(nil, mongo.ErrNoDocuments)
	mgs := svc.DeleteRating(ratingIdFailed)
	assert.Equal(t, message.ErrDataNotFound, mgs)
}

func TestDeleteRatingFailed4(t *testing.T) {
	rId := "629ec07e6f3c2761ba2dc419"
	ObjRatingId, _ := primitive.ObjectIDFromHex(rId)

	rating := &entity.RatingsCol{
		ID:   ObjRatingId,
		Name: name + "22",
	}

	ratingSubmission := entity.RatingSubmisson{
		RatingID: rId,
	}

	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(rating, nil)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", rId).Return(ratingSubmission)
	mgs := svc.DeleteRating(rId)
	assert.Equal(t, message.ErrRatingHasRatingSubmission, mgs)
}

func TestDeleteRatingFailed5(t *testing.T) {
	ObjRatingId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")

	rating := &entity.RatingsCol{
		ID:   ObjRatingId,
		Name: name + "23",
	}

	ratingSubmission := entity.RatingSubmisson{
		RatingID: ratingId,
	}

	ratingRepository.Mock.On("GetRatingById", ObjRatingId).Return(rating, nil)
	ratingRepository.Mock.On("GetRatingSubmissionByRatingId", "629dce7bf1f26275e0d84826").Return(ratingSubmission)
	mgs := svc.DeleteRating("629dce7bf1f26275e0d84826")
	assert.Equal(t, message.FailedMsg, mgs)
}

func TestGetListRatings(t *testing.T) {
	req := request.GetListRatingsRequest{
		Sort:  "",
		Dir:   "desc",
		Page:  0,
		Limit: 0,
	}
	objectId1, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc468")
	objectId2, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc848")
	result := []entity.RatingsCol{
		{
			ID:          objectId1,
			Description: &description,
		},
		{
			ID:          objectId2,
			Description: &description,
		},
	}
	paginationResult := base.Pagination{
		Records:   2,
		Limit:     50,
		Page:      1,
		TotalPage: 1,
	}
	ratingRepository.Mock.On("GetRatingsByParams", request.RatingFilter{SourceUid: []string(nil), RatingTypeId: []string(nil)}, 1, 50, "updated_at", -1).Return(result, &paginationResult, nil)

	_, _, msg := svc.GetListRatings(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetListRatingSummary(t *testing.T) {
	req := request.GetListRatingSummaryRequest{
		Sort:   "",
		Dir:    "desc",
		Page:   0,
		Limit:  0,
		Filter: "{\"source_uid\": [\"2729\", \"2951\"],\"score\":[0,5]}",
	}
	result := []entity.RatingSubmisson{
		{
			RatingID: "629ec07e6f3c2761ba2dc468",
			Comment:  &Desc,
			Value:    value,
		},
		{
			RatingID: "629ec07e6f3c2761ba2dc848",
			Comment:  &Desc,
			Value:    value,
		},
	}
	objectId1, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc468")
	objectId2, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc848")
	result2 := []entity.RatingsCol{
		{
			ID:          objectId1,
			Description: &description,
			SourceUid:   "2951",
		},
		{
			ID:          objectId2,
			Description: &description,
			SourceUid:   "2729",
		},
	}
	paginationResult := base.Pagination{
		Records:   2,
		Limit:     50,
		Page:      1,
		TotalPage: 1,
	}
	ratingRepository.Mock.On("GetRatingsByParams", request.RatingFilter{SourceUid: []string{"2729", "2951"}, RatingTypeId: []string(nil)}, 1, 50, "updated_at", -1).Return(result2, &paginationResult, nil)
	ratingRepository.Mock.On("GetListRatingSubmissions", request.RatingSubmissionFilter{UserIDLegacy: []string(nil), Score: []float64{0, 5}, RatingID: []string{"629ec07e6f3c2761ba2dc468", "629ec07e6f3c2761ba2dc848"}, StartDate: "", EndDate: ""}, 1, int64(50), "updated_at", -1).Return(result, &paginationResult, nil)
	_, msg := svc.GetListRatingSummary(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetListRatingSubmissionErrWrongFilter(t *testing.T) {
	input := request.ListRatingSubmissionRequest{
		Dir:    "desc",
		Sort:   "wrong filter",
		Filter: "{\"user_uid\":[\"a12346fb-bd93-fedc-abcd-0739865540cb\",\"0739865540cb-bd93-fedc-abcd-a12346fb\"],\"score\":[4]}",
	}

	filter := request.RatingSubmissionFilter{}
	_ = json.Unmarshal([]byte(input.Filter), &filter)

	page := base.Pagination{
		Records:      0,
		TotalRecords: 0,
		Limit:        50,
		Page:         1,
	}

	ratingRepository.Mock.On("GetListRatingSubmissions", filter, 1, int64(page.Limit), input.Sort, -1).Return(nil, &page, gorm.ErrInvalidDB)

	_, _, msg := svc.GetListRatingSubmissions(input)

	assert.Equal(t, message.FailedMsg, msg)
}

func TestUpdateRatingSubmissionErrMarshall1(t *testing.T) {
	minScore := 0
	maxScore := 5
	intervals := 6
	valueRate := "1"
	numericID := "629dce7bf1f26275e0d84826"
	likertID := "629dce7bf1f26275e0d84826"
	objectId, _ := primitive.ObjectIDFromHex(id)
	input := request.UpdateRatingSubmissionRequest{
		ID:       "id",
		Value:    &valueRate,
		RatingID: id,
	}

	sub := entity.RatingSubmisson{
		UserID:       &id,
		UserIDLegacy: &id,
		RatingID:     id,
	}

	rating := entity.RatingsCol{
		ID:             objectId,
		RatingTypeId:   id,
		CommentAllowed: &Bool,
		Status:         &Bool,
	}

	num := entity.RatingTypesNumCol{
		ID:          objectId,
		Status:      &Bool,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Description: &Desc,
		Scale:       &Scale,
		Intervals:   &intervals,
	}

	objectNumericId, _ := primitive.ObjectIDFromHex(numericID)
	objectLikertId, _ := primitive.ObjectIDFromHex(likertID)
	ratingRepository.Mock.On("GetRatingSubmissionById", objectId).Return(sub, nil)
	ratingRepository.Mock.On("FindRatingByRatingID", objectId).Return(rating, nil)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDLegacyAndRatingID", &id, id, id).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", &id, id, id).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingNumericTypeByRatingTypeID", objectNumericId).Return(num, nil)
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", objectLikertId).Return(num, nil)
	ratingRepository.Mock.On("UpdateRatingSubmission", input).Return(sub, nil)

	msg := svc.UpdateRatingSubmission(input)

	assert.Equal(t, message.RatingSubmissionNotFound, msg)
}

func TestUpdateRatingSubmissionErrMarshall2(t *testing.T) {
	minScore := 0
	maxScore := 5
	intervals := 6
	valueRate := "1"
	numericID := "629dce7bf1f26275e0d84826"
	likertID := "629dce7bf1f26275e0d84826"
	objectId, _ := primitive.ObjectIDFromHex(id)
	input := request.UpdateRatingSubmissionRequest{
		ID:       id,
		Value:    &valueRate,
		RatingID: "id",
	}

	sub := entity.RatingSubmisson{
		UserID:       &id,
		UserIDLegacy: &id,
		RatingID:     id,
	}

	rating := entity.RatingsCol{
		ID:             objectId,
		RatingTypeId:   id,
		CommentAllowed: &Bool,
		Status:         &Bool,
	}

	num := entity.RatingTypesNumCol{
		ID:          objectId,
		Status:      &Bool,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Description: &Desc,
		Scale:       &Scale,
		Intervals:   &intervals,
	}

	objectNumericId, _ := primitive.ObjectIDFromHex(numericID)
	objectLikertId, _ := primitive.ObjectIDFromHex(likertID)
	ratingRepository.Mock.On("GetRatingSubmissionById", objectId).Return(sub, nil)
	ratingRepository.Mock.On("FindRatingByRatingID", objectId).Return(rating, nil)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDLegacyAndRatingID", &id, id, id).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", &id, id, id).Return(nil, gorm.ErrRecordNotFound)
	ratingRepository.Mock.On("FindRatingNumericTypeByRatingTypeID", objectNumericId).Return(num, nil)
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", objectLikertId).Return(num, nil)
	ratingRepository.Mock.On("UpdateRatingSubmission", input).Return(sub, nil)

	msg := svc.UpdateRatingSubmission(input)

	assert.Equal(t, message.ErrRatingNotFound, msg)
}

func TestGetListRatingSummaryWrongFilter(t *testing.T) {
	req := request.GetListRatingSummaryRequest{
		Sort:   "wrong filter",
		Dir:    "desc",
		Page:   0,
		Limit:  0,
		Filter: "{\"source_uid\": [\"2729\", \"2951\"]}",
	}

	objectId1, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc468")
	objectId2, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc848")
	result2 := []entity.RatingsCol{
		{
			ID:          objectId1,
			Description: &description,
			SourceUid:   "2951",
		},
		{
			ID:          objectId2,
			Description: &description,
			SourceUid:   "2729",
		},
	}
	paginationResult := base.Pagination{
		Records:   2,
		Limit:     50,
		Page:      1,
		TotalPage: 1,
	}
	ratingRepository.Mock.On("GetRatingsByParams", request.RatingFilter{SourceUid: []string{"2729", "2951"}, RatingTypeId: []string(nil)}, 1, 50, "wrong filter", -1).Return(result2, &paginationResult, nil)
	ratingRepository.Mock.On("GetListRatingSubmissions", request.RatingSubmissionFilter{UserIDLegacy: []string(nil), Score: []float64(nil), RatingID: []string{"629ec07e6f3c2761ba2dc468", "629ec07e6f3c2761ba2dc848"}, StartDate: "", EndDate: ""}, 1, int64(50), "wrong filter", -1).Return(nil, &paginationResult, gorm.ErrInvalidDB)
	_, msg := svc.GetListRatingSummary(req)
	assert.Equal(t, message.WrongFilter, msg)
}

func TestGetListRatingSummaryWrongScoreFilter(t *testing.T) {
	req := request.GetListRatingSummaryRequest{
		Sort:   "",
		Dir:    "desc",
		Page:   0,
		Limit:  0,
		Filter: "{\"source_uid\": [\"2729\", \"2951\"], \"score\":[4,5,6]}",
	}
	result := []entity.RatingSubmisson{
		{
			RatingID: "629ec07e6f3c2761ba2dc468",
			Comment:  &Desc,
			Value:    value,
		},
		{
			RatingID: "629ec07e6f3c2761ba2dc848",
			Comment:  &Desc,
			Value:    value,
		},
	}
	objectId1, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc468")
	objectId2, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc848")
	result2 := []entity.RatingsCol{
		{
			ID:          objectId1,
			Description: &description,
			SourceUid:   "2951",
		},
		{
			ID:          objectId2,
			Description: &description,
			SourceUid:   "2729",
		},
	}
	paginationResult := base.Pagination{
		Records:   2,
		Limit:     50,
		Page:      1,
		TotalPage: 1,
	}
	ratingRepository.Mock.On("GetRatingsByParams", request.RatingFilter{SourceUid: []string{"2729", "2951"}, RatingTypeId: []string(nil)}, 1, 50, "updated_at", -1).Return(result2, &paginationResult, nil)
	ratingRepository.Mock.On("GetListRatingSubmissions", request.RatingSubmissionFilter{UserIDLegacy: []string(nil), Score: []float64{4, 5, 6}, RatingID: []string{"629ec07e6f3c2761ba2dc468", "629ec07e6f3c2761ba2dc848"}, StartDate: "", EndDate: ""}, 1, int64(50), "updated_at", -1).Return(result, &paginationResult, nil)
	_, msg := svc.GetListRatingSummary(req)
	assert.Equal(t, message.WrongScoreFilter, msg)
}

func TestGetListRatingSummaryNoScoreFilter(t *testing.T) {
	req := request.GetListRatingSummaryRequest{
		Sort:   "",
		Dir:    "asc",
		Page:   0,
		Limit:  0,
		Filter: "{\"source_uid\": [\"2729\", \"2951\"]}",
	}
	result := []entity.RatingSubmisson{
		{
			RatingID: "629ec07e6f3c2761ba2dc468",
			Comment:  &Desc,
			Value:    value,
		},
		{
			RatingID: "629ec07e6f3c2761ba2dc848",
			Comment:  &Desc,
			Value:    value,
		},
	}
	objectId1, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc468")
	objectId2, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc848")
	result2 := []entity.RatingsCol{
		{
			ID:          objectId1,
			Description: &description,
			SourceUid:   "2951",
		},
		{
			ID:          objectId2,
			Description: &description,
			SourceUid:   "2729",
		},
	}
	paginationResult := base.Pagination{
		Records:   2,
		Limit:     50,
		Page:      1,
		TotalPage: 1,
	}
	ratingRepository.Mock.On("GetRatingsByParams", request.RatingFilter{SourceUid: []string{"2729", "2951"}, RatingTypeId: []string(nil)}, 1, 50, "updated_at", 1).Return(result2, &paginationResult, nil)
	ratingRepository.Mock.On("GetListRatingSubmissions", request.RatingSubmissionFilter{UserIDLegacy: []string(nil), Score: []float64(nil), RatingID: []string{"629ec07e6f3c2761ba2dc468", "629ec07e6f3c2761ba2dc848"}, StartDate: "", EndDate: ""}, 1, int64(50), "updated_at", 1).Return(result, &paginationResult, nil)
	_, msg := svc.GetListRatingSummary(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetListRatingSummaryErrSourceUidRequired(t *testing.T) {
	req := request.GetListRatingSummaryRequest{
		Sort:   "",
		Dir:    "desc",
		Page:   0,
		Limit:  0,
		Filter: "",
	}
	result := []entity.RatingSubmisson{
		{
			RatingID: "629ec07e6f3c2761ba2dc468",
			Comment:  &Desc,
			Value:    value,
		},
		{
			RatingID: "629ec07e6f3c2761ba2dc848",
			Comment:  &Desc,
			Value:    value,
		},
	}
	objectId1, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc468")
	objectId2, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc848")
	result2 := []entity.RatingsCol{
		{
			ID:          objectId1,
			Description: &description,
			SourceUid:   "2951",
		},
		{
			ID:          objectId2,
			Description: &description,
			SourceUid:   "2729",
		},
	}
	paginationResult := base.Pagination{
		Records:   2,
		Limit:     50,
		Page:      1,
		TotalPage: 1,
	}
	ratingRepository.Mock.On("GetRatingsByParams", request.RatingFilter{SourceUid: []string{"2729", "2951"}, RatingTypeId: []string(nil)}, 1, 50, "updated_at", -1).Return(result2, &paginationResult, nil)
	ratingRepository.Mock.On("GetListRatingSubmissions", request.RatingSubmissionFilter{UserIDLegacy: []string(nil), Score: []float64(nil), RatingID: []string{"629ec07e6f3c2761ba2dc468", "629ec07e6f3c2761ba2dc848"}, StartDate: "", EndDate: ""}, 1, int64(50), "updated_at", -1).Return(result, &paginationResult, nil)
	_, msg := svc.GetListRatingSummary(req)
	assert.Equal(t, message.ErrSourceUidRequire, msg)
}

func TestCancelRatingSubmissionSuccess(t *testing.T) {
	ids := []primitive.ObjectID{}
	input := request.CancelRatingById{
		RatingSubmissionId: []string{"630dca3fc27e5483bdc006ec", "630dca3fc27e5483bdc006ed"},
		CancelledReason:    "Cancelled Reason Test",
	}

	for _, id := range input.RatingSubmissionId {
		objectId, _ := primitive.ObjectIDFromHex(id)
		ids = append(ids, objectId)
	}

	ratingRepository.Mock.On("CancelRatingSubmissionByIds", ids, input.CancelledReason).Return(nil)
	msg := svc.CancelRatingSubmission(input)

	assert.Equal(t, message.SuccessMsg, msg)
}

func TestCancelRatingSubmissionFailed(t *testing.T) {
	ids := []primitive.ObjectID{}
	input := request.CancelRatingById{
		RatingSubmissionId: []string{"630dca3fc27e5483bdc006ec", "630dca3fc27e5483bdc006ed"},
		CancelledReason:    "failed",
	}

	for _, id := range input.RatingSubmissionId {
		objectId, _ := primitive.ObjectIDFromHex(id)
		ids = append(ids, objectId)
	}

	ratingRepository.Mock.On("CancelRatingSubmissionByIds", ids, input.CancelledReason).Return(errors.New("failed"))
	msg := svc.CancelRatingSubmission(input)

	assert.Equal(t, message.ErrSaveData, msg)
}
