package test

import (
	"encoding/json"
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/message"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger log.Logger

var ratingRepository = &repository_mock.RatingRepositoryMock{Mock: mock.Mock{}}
var svc = service.NewRatingService(logger, ratingRepository)

func init() {
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}
}

var (
	id          = "629dce7bf1f26275e0d84826"
	description = "alo"
	valueFailed = "failed"

	Desc          = "Description"
	Bool          = true
	Scale         = 0
	value float64 = 4
)

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
		Type:        "12345",
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
		Type:        "12345",
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
		Type:        "12345",
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
		Type:        "12345",
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
		Comment:  "dlfslkjdf",
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
		Type:        "12345",
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
		Type:        "12345",
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
		Comment:  "dlfslkjdf",
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
		ID: objectId,
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

	assert.Equal(t, message.ErrNoData, msg)
}

func TestGetRatingSubmissionWrongID(t *testing.T) {
	failId := "123456"
	objectId, _ := primitive.ObjectIDFromHex("123456")

	ratingRepository.Mock.On("GetRatingSubmissionById", objectId).Return(nil, mongo.ErrNoDocuments)

	_, msg := svc.GetRatingSubmission(failId)

	assert.Equal(t, message.ErrDataNotFound, msg)
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

func TestCreateRatingSubmissionSuccess(t *testing.T) {
	minScore := 0
	maxScore := 5
	intervals := 6
	matchStrValuePtr := "match"
	objectId, _ := primitive.ObjectIDFromHex(id)
	input := request.CreateRatingSubmissonRequest{
		UserID:       &matchStrValuePtr,
		UserIDLegacy: &matchStrValuePtr,
		RatingID:     id,
		Value:        &value,
		UserAgent:    "user agent",
	}

	sub := entity.RatingSubmisson{
		UserID:       &matchStrValuePtr,
		UserIDLegacy: &matchStrValuePtr,
		RatingID:     id,
	}

	rating := entity.RatingsCol{
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

	ratingRepository.Mock.On("FindRatingByRatingID", objectId).Return(rating, nil)
	ratingRepository.Mock.On("FindRatingNumericTypeByRatingTypeID", objectId).Return(num, nil)
	//ratingRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", &matchStrValuePtr, id).Return(nil, errors.New("record found"))
	ratingRepository.Mock.On("CreateRatingSubmission", input).Return(sub, nil)

	msg := svc.CreateRatingSubmission(input)

	assert.Equal(t, message.SuccessMsg, msg)
}

func TestCreateRatingSubmissionWrongValue(t *testing.T) {
	minScore := 0
	maxScore := 5
	intervals := 6
	var vl float64 = 4.5
	matchStrValuePtr := "match"
	objectId, _ := primitive.ObjectIDFromHex(id)
	input := request.CreateRatingSubmissonRequest{
		UserID:       &matchStrValuePtr,
		UserIDLegacy: &matchStrValuePtr,
		RatingID:     id,
		Value:        &vl,
		UserAgent:    "user agent",
	}

	sub := entity.RatingSubmisson{
		UserID:       &matchStrValuePtr,
		UserIDLegacy: &matchStrValuePtr,
		RatingID:     id,
	}

	rating := entity.RatingsCol{
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

	ratingRepository.Mock.On("FindRatingByRatingID", objectId).Return(rating, nil)
	ratingRepository.Mock.On("FindRatingNumericTypeByRatingTypeID", objectId).Return(num, nil)
	//ratingRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", &matchStrValuePtr, id).Return(nil, errors.New("record found"))
	ratingRepository.Mock.On("CreateRatingSubmission", input).Return(sub, nil)

	msg := svc.CreateRatingSubmission(input)

	assert.Equal(t, message.ErrValueFormat.Code, msg.Code)
}

func TestCreateRatingSubmissionRequireID(t *testing.T) {
	minScore := 0
	maxScore := 5
	intervals := 6
	matchStrValuePtr := "match"
	var vl float64 = 4.5
	objectId, _ := primitive.ObjectIDFromHex(id)
	input := request.CreateRatingSubmissonRequest{
		UserID:       nil,
		UserIDLegacy: nil,
		RatingID:     id,
		Value:        &vl,
		UserAgent:    "user agent",
	}

	sub := entity.RatingSubmisson{
		UserID:       &matchStrValuePtr,
		UserIDLegacy: &matchStrValuePtr,
		RatingID:     id,
	}

	rating := entity.RatingsCol{
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

	ratingRepository.Mock.On("FindRatingByRatingID", objectId).Return(rating, nil)
	ratingRepository.Mock.On("FindRatingNumericTypeByRatingTypeID", objectId).Return(num, nil)
	//ratingRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", &matchStrValuePtr, id).Return(nil, errors.New("record found"))
	ratingRepository.Mock.On("CreateRatingSubmission", input).Return(sub, nil)

	msg := svc.CreateRatingSubmission(input)

	assert.Equal(t, message.UserUIDRequired, msg)
}

func TestUpdateRatingSubmissionSuccess(t *testing.T) {
	minScore := 0
	maxScore := 5
	intervals := 6
	matchStrValuePtr := "match"
	objectId, _ := primitive.ObjectIDFromHex(id)
	input := request.UpdateRatingSubmissonRequest{
		UserID:       &matchStrValuePtr,
		UserIDLegacy: &matchStrValuePtr,
		RatingID:     id,
		Value:        4,
	}

	sub := entity.RatingSubmisson{
		UserID:       &matchStrValuePtr,
		UserIDLegacy: &matchStrValuePtr,
		RatingID:     id,
	}

	rating := entity.RatingsCol{
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

	ratingRepository.Mock.On("FindRatingByRatingID", objectId).Return(rating, nil)
	ratingRepository.Mock.On("FindRatingNumericTypeByRatingTypeID", objectId).Return(num, nil)
	//ratingRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", &matchStrValuePtr, id).Return(nil, errors.New("record found"))
	ratingRepository.Mock.On("CreateRatingSubmission", input).Return(sub, nil)

	msg := svc.UpdateRatingSubmission(input)

	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetListRatingSubmission(t *testing.T) {
	matchStrValuePtr := "match"
	input := request.ListRatingSubmissionRequest{
		Dir:    "asc",
		Filter: "{\"user_uid\":[\"a12346fb-bd93-fedc-abcd-0739865540cb\",\"0739865540cb-bd93-fedc-abcd-a12346fb\"],\"score\":[4,4.5]}",
	}

	filter := request.RatingSubmissionFilter{}
	_ = json.Unmarshal([]byte(input.Filter), &filter)
	subs := []entity.RatingSubmisson{
		{
			UserID:       &matchStrValuePtr,
			UserIDLegacy: &matchStrValuePtr,
			Value:        4.5,
		},
	}

	page := base.Pagination{
		Records:      1,
		TotalRecords: 1,
		Limit:        50,
		Page:         1,
	}

	ratingRepository.Mock.On("GetListRatingSubmissions", filter, 1, int64(page.Limit), "updated_at", 1).Return(subs, &page, nil)

	_, _, msg := svc.GetListRatingSubmissions(input)

	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetListRatingSubmissionMarshalErr(t *testing.T) {
	matchStrValuePtr := "match"
	input := request.ListRatingSubmissionRequest{
		Dir:    "asc",
		Filter: "{\"user_uid\":[\"a12346fb-bd93-fedc-abcd-0739865540cb\",\"0739865540cb-bd93-fedc-abcd-a12346fb\"],\"score\":[\"4\",\"4.5\"]}",
	}

	filter := request.RatingSubmissionFilter{}
	_ = json.Unmarshal([]byte(input.Filter), &filter)
	subs := []entity.RatingSubmisson{
		{
			UserID:       &matchStrValuePtr,
			UserIDLegacy: &matchStrValuePtr,
			Value:        4.5,
		},
	}

	page := base.Pagination{
		Records:      1,
		TotalRecords: 1,
		Limit:        50,
		Page:         1,
	}

	ratingRepository.Mock.On("GetListRatingSubmissions", filter, 1, int64(page.Limit), "updated_at", 1).Return(subs, &page, nil)

	_, _, msg := svc.GetListRatingSubmissions(input)

	assert.Equal(t, message.ErrUnmarshalFilterListRatingRequest, msg)
}

func TestCreateRatingTypeLikert(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Type:        "type",
		Description: &description,
	}

	msg := svc.CreateRatingTypeLikert(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestCreateRatingTypeLikertFailed(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Type:        "typeErr",
		Description: &description,
	}

	msg := svc.CreateRatingTypeLikert(req)
	assert.Equal(t, message.FailedMsg, msg)
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

func TestUpdateRatingTypeLikert(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:          "629ec07e6f3c2761ba2dc868",
		Description: &description,
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)

	likert := entity.RatingTypesLikertCol{
		Type:        "test",
		Description: &description,
	}
	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(likert)

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestUpdateRatingTypeLikertErrIdFormatReq(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:          "213",
		Description: &description,
	}

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.ErrNoData, msg)
}

func TestUpdateRatingTypeLikertErrSaveData(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:            "629ec07e6f3c2761ba2dc828",
		NumStatements: 1,
		Statement01:   &description,
		Statement02:   &description,
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)

	likert := entity.RatingTypesLikertCol{
		Type:          "test",
		NumStatements: 2,
	}
	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(likert)

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.ErrMatchNumState, msg)
}

func TestUpdateRatingTypeLikertErrSaveData2(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:            "629ec07e6f3c2761ba2dc828",
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

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.ErrMatchNumState, msg)
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
		Comment:  "dlfslkjdf",
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
