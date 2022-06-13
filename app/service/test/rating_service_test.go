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
	id = "629dce7bf1f26275e0d84826"
)

func TestCreateRatingTypeNum(t *testing.T) {
	var minScore = 0
	var scale = 1
	var status bool
	req := request.CreateRatingTypeNumRequest{
		Type:        "type",
		Description: "Description",
		MinScore:    &minScore,
		MaxScore:    1,
		Scale:       &scale,
		Intervals:   11,
		Status:      &status,
	}
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	ratingTypesNumCol := entity.RatingTypesNumCol{
		ID: objectId,
	}

	ratingRepository.Mock.On("CreateRatingTypeNum", req).Return(ratingTypesNumCol)

	msg := svc.CreateRatingTypeNum(req)
	assert.Equal(t, message.SuccessMsg.Message, msg.Message)
}

func TestCreateRatingTypeNumErrScaleValueReq(t *testing.T) {
	var minScore = 0
	var scale = 3
	var status bool
	req := request.CreateRatingTypeNumRequest{
		Type:        "type",
		Description: "Description",
		MinScore:    &minScore,
		MaxScore:    2,
		Scale:       &scale,
		Intervals:   11,
		Status:      &status,
	}
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	ratingTypesNumCol := entity.RatingTypesNumCol{
		ID: objectId,
	}

	ratingRepository.Mock.On("CreateRatingTypeNum", req).Return(ratingTypesNumCol)

	msg := svc.CreateRatingTypeNum(req)
	assert.Equal(t, message.ErrScaleValueReq, msg)
}

func TestCreateRatingTypeNumErrSaveData(t *testing.T) {
	var minScore = 0
	var scale = 1
	var status bool
	req := request.CreateRatingTypeNumRequest{
		Type:        "12345",
		Description: "Description",
		MinScore:    &minScore,
		MaxScore:    1,
		Scale:       &scale,
		Intervals:   11,
		Status:      &status,
	}
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	ratingTypesNumCol := entity.RatingTypesNumCol{
		ID: objectId,
	}

	ratingRepository.Mock.On("CreateRatingTypeNum", req).Return(ratingTypesNumCol)

	msg := svc.CreateRatingTypeNum(req)
	assert.Equal(t, message.ErrSaveData, msg)
}

func TestGetRatingTypeNumById(t *testing.T) {
	req := request.GetRatingTypeNumRequest{Id: "629ec0736f3c2761ba2dc867"}
	var minScore = 0
	var scale = 2
	var status bool

	objectId, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")
	ratingTypesNumCol := entity.RatingTypesNumCol{
		ID:          objectId,
		Type:        "type",
		Description: "Description",
		MinScore:    &minScore,
		MaxScore:    2,
		Scale:       &scale,
		Intervals:   11,
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

	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	ratingTypesNumCol := entity.RatingTypesNumCol{
		ID:          objectId,
		Type:        "type",
		Description: "Description",
		MinScore:    &minScore,
		MaxScore:    2,
		Scale:       &scale,
		Intervals:   11,
		Status:      &status,
	}

	ratingRepository.Mock.On("GetRatingTypeNumById", objectId).Return(ratingTypesNumCol)

	_, msg := svc.GetRatingTypeNumById(req)
	assert.Equal(t, message.FailedMsg, msg)
}

func TestGetRatingTypeNumByIdErrNoData(t *testing.T) {
	req := request.GetRatingTypeNumRequest{Id: "629ec07e6f3c2761ba2dc868"}

	objectId, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc868")

	ratingRepository.Mock.On("GetRatingTypeNumById", objectId).Return(nil)

	_, msg := svc.GetRatingTypeNumById(req)
	assert.Equal(t, message.ErrNoData, msg)
}

func TestUpdateRatingTypeNum(t *testing.T) {
	var minScore = 0
	var scale = 2
	var status bool
	req := request.CreateRatingTypeNumRequest{
		Id:          "629ec0736f3c2761ba2dc867",
		Type:        "12345",
		Description: "Description",
		MinScore:    &minScore,
		MaxScore:    1,
		Scale:       &scale,
		Intervals:   11,
		Status:      &status,
	}
	objectId, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")

	ratingRepository.Mock.On("UpdateRatingTypeNum", objectId)

	msg := svc.UpdateRatingTypeNum(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestUpdateRatingTypeNumFailed(t *testing.T) {
	var minScore = 0
	var scale = 2
	var status bool
	req := request.CreateRatingTypeNumRequest{
		Id:          "629ec07e6f3c2761ba2dc868",
		Type:        "12345",
		Description: "Description",
		MinScore:    &minScore,
		MaxScore:    1,
		Scale:       &scale,
		Intervals:   11,
		Status:      &status,
	}
	objectId, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc868")

	ratingRepository.Mock.On("UpdateRatingTypeNum", objectId)

	msg := svc.UpdateRatingTypeNum(req)
	assert.Equal(t, message.ErrSaveData, msg)
}

func TestUpdateRatingTypeNumErrIdFormatReq(t *testing.T) {
	var minScore = 0
	var scale = 2
	var status bool
	req := request.CreateRatingTypeNumRequest{
		Id:          "sdkj234kld",
		Type:        "12345",
		Description: "Description",
		MinScore:    &minScore,
		MaxScore:    1,
		Scale:       &scale,
		Intervals:   11,
		Status:      &status,
	}
	objectId, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc868")

	ratingRepository.Mock.On("UpdateRatingTypeNum", objectId)

	msg := svc.UpdateRatingTypeNum(req)
	assert.Equal(t, message.ErrIdFormatReq, msg)
}

func TestDeleteRatingTypeNumById(t *testing.T) {
	req := request.GetRatingTypeNumRequest{Id: "629ec0736f3c2761ba2dc867"}

	msg := svc.DeleteRatingTypeNumById(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestDeleteRatingTypeNumByIdErrIdFormatReq(t *testing.T) {
	req := request.GetRatingTypeNumRequest{Id: "rj2j3lk4324"}

	msg := svc.DeleteRatingTypeNumById(req)
	assert.Equal(t, message.ErrIdFormatReq, msg)
}

func TestDeleteRatingTypeNumByIdFailed(t *testing.T) {
	req := request.GetRatingTypeNumRequest{Id: "629ec0836f3c2761ba2dc899"}

	msg := svc.DeleteRatingTypeNumById(req)
	assert.Equal(t, message.FailedMsg, msg)
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
			Description: "jdhkaf",
		},
		{
			ID:          objectId2,
			Description: "jdhkaf",
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

func TestCreateRatingSubmissionSuccess(t *testing.T) {
	matchStrValuePtr := "match"
	objectId, _ := primitive.ObjectIDFromHex(id)
	input := request.CreateRatingSubmissonRequest{
		UserID:       &matchStrValuePtr,
		UserIDLegacy: &matchStrValuePtr,
		RatingID:     id,
		Value:        4.5,
		UserAgent:    "user agent",
	}

	sub := entity.RatingSubmisson{
		UserID:   &matchStrValuePtr,
		RatingID: id,
	}

	rating := entity.RatingsCol{
		RatingTypeId:   id,
		CommentAllowed: nil,
		Status:         nil,
	}

	num := entity.RatingTypesNumCol{
		ID:     objectId,
		Status: nil,
	}

	ratingRepository.Mock.On("FindRatingByRatingID", objectId).Return(rating, nil)
	ratingRepository.Mock.On("FindRatingNumericTypeByRatingTypeID", objectId).Return(num, nil)
	//ratingRepository.Mock.On("FindRatingSubmissionByUserIDAndRatingID", &matchStrValuePtr, id).Return(nil, errors.New("record found"))
	ratingRepository.Mock.On("CreateRatingSubmission", input).Return(sub, nil)

	msg := svc.CreateRatingSubmission(input)

	assert.Equal(t, message.SuccessMsg, msg)
}

func TestUpdateRatingSubmissionSuccess(t *testing.T) {
	matchStrValuePtr := "match"
	objectId, _ := primitive.ObjectIDFromHex(id)
	input := request.UpdateRatingSubmissonRequest{
		UserID:       &matchStrValuePtr,
		UserIDLegacy: &matchStrValuePtr,
		RatingID:     id,
		Value:        4.5,
	}

	sub := entity.RatingSubmisson{
		UserID:   &matchStrValuePtr,
		RatingID: id,
	}

	rating := entity.RatingsCol{
		RatingTypeId:   id,
		CommentAllowed: nil,
		Status:         nil,
	}

	num := entity.RatingTypesNumCol{
		ID:     objectId,
		Status: nil,
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

	assert.Equal(t, message.SuccessMsg, msg)
}

func TestCreateRatingTypeLikert(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Type:        "type",
		Description: "Description",
	}

	msg := svc.CreateRatingTypeLikert(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestCreateRatingTypeLikertFailed(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Type:        "typeErr",
		Description: "Description",
	}

	msg := svc.CreateRatingTypeLikert(req)
	assert.Equal(t, message.ErrSaveData, msg)
}

func TestGetRatingTypeLikertById(t *testing.T) {
	req := request.GetRatingTypeLikertRequest{
		Id: "629ec07e6f3c2761ba2dc868",
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)

	likert := entity.RatingTypesLikertCol{
		Type:        "test",
		Description: "dkfjlsdf",
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
	assert.Equal(t, message.ErrIdFormatReq, msg)
}

func TestGetRatingTypeLikertByIdFailed(t *testing.T) {
	req := request.GetRatingTypeLikertRequest{
		Id: "629dce7bf1f26275e0d84826",
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)

	likert := entity.RatingTypesLikertCol{
		Type:        "test",
		Description: "dkfjlsdf",
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
		Description: "fjkdsfd",
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)

	likert := entity.RatingTypesLikertCol{
		Type:        "test",
		Description: "dkfjlsdf",
	}
	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(likert)

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestUpdateRatingTypeLikertErrIdFormatReq(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:          "213",
		Description: "fjkdsfd",
	}

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.ErrIdFormatReq, msg)
}

func TestUpdateRatingTypeLikertErrSaveData(t *testing.T) {
	req := request.SaveRatingTypeLikertRequest{
		Id:          "629ec07e6f3c2761ba2dc828",
		Description: "failed",
	}
	objectId, _ := primitive.ObjectIDFromHex(req.Id)

	likert := entity.RatingTypesLikertCol{
		Type:        "test",
		Description: "dkfjlsdf",
	}
	ratingRepository.Mock.On("GetRatingTypeLikertById", objectId).Return(likert)

	msg := svc.UpdateRatingTypeLikert(req)
	assert.Equal(t, message.ErrSaveData, msg)
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
	assert.Equal(t, message.ErrIdFormatReq, msg)
}

func TestDeleteRatingTypeLikertByIdFailedMsg(t *testing.T) {
	req := request.GetRatingTypeLikertRequest{
		Id: "629ec0836f3c2761ba2dc899",
	}
	msg := svc.DeleteRatingTypeLikertById(req)
	assert.Equal(t, message.FailedMsg, msg)
}

func TestDeleteRatingTypeLikertByIdErrNoData(t *testing.T) {
	req := request.GetRatingTypeLikertRequest{
		Id: "629ec0836f3c2761ba2dc869",
	}
	msg := svc.DeleteRatingTypeLikertById(req)
	assert.Equal(t, message.ErrNoData, msg)
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
			Description: "jdhkaf",
		},
		{
			ID:          objectId2,
			Description: "jdhkaf",
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
			Description: "jdhkaf",
		},
		{
			ID:          objectId2,
			Description: "jdhkaf",
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
