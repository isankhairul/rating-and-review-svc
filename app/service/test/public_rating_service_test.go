package test

import (
	"errors"
	"go-klikdokter/app/model/base"
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
	"go.mongodb.org/mongo-driver/mongo"
)

var publicRatingRepository = &repository_mock.PublicRatingRepositoryMock{Mock: mock.Mock{}}
var publicRatingService = service.NewPublicRatingService(logger, ratingRepository, publicRatingRepository)

func init() {
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}
}

var (
	userId            = "2210"
	comment           = "Comment Test"
	ipaddress         = "138.199.20.50"
	useragent         = "Chrome/{Chrome Rev} Mobile Safari/{WebKit Rev}"
	sourceTransId     = "2123"
	ratingid          = "62c4f30f6d90d90d6594fab9"
	ratingType        = "performance-doctor"
	ratingSubId       = "629dce7bf1f26275e0d84826"
	ratingSubHelpId   = "62c6438c08d23eb8fe9834e8"
	ratingSubIdFailed = "62c53baf039c7a6554accb0d"
	idDummy1          = "62c4f2b96d90d90d6594fab7"
	idDummy2          = "62c4f30f6d90d90d6594fab8"
	failedId          = "62c3e57b457ed515928c3690"
	displayName       = "User Name"
	anonym            = false
)

var requestSummary = request.GetPublicListRatingSummaryRequest{
	SourceType: "doctor",
	Sort:       "",
	Dir:        "desc",
	Page:       1,
	Limit:      50,
}

var requestSubmission = request.GetPublicListRatingSubmissionRequest{
	SourceType: "doctor",
	SourceUID:  "895",
	Sort:       "created_at",
	Dir:        "desc",
	Page:       1,
	Limit:      50,
}

var filterSummary = request.FilterRatingSummary{
	SourceType: requestSummary.SourceType,
	SourceUid:  []string(nil),
	RatingType: []string(nil),
}

func TestGetRatingBySourceTypeAndSourceUID(t *testing.T) {
	req := request.GetRatingBySourceTypeAndActorRequest{
		SourceType: "doctor",
		SourceUID:  "894",
	}

	ratingId, _ := primitive.ObjectIDFromHex("629ec0736f3c2761ba2dc867")
	ratingTypeId, _ := primitive.ObjectIDFromHex("62c4f03b6d90d90d6594fab5")
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
	filter := request.GetRatingBySourceTypeAndActorFilter{}

	publicRatingRepository.Mock.On("GetRatingsBySourceTypeAndActor", req.SourceType, req.SourceUID, filter).Return(resultRatings, nil).Once()
	publicRatingRepository.Mock.On("GetRatingTypeLikertById", ratingTypeId).Return(resultRatingType, nil).Once()

	_, msg := svc.GetRatingBySourceTypeAndActor(req)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetRatingBySourceTypeAndSourceUIDErrNoDataRating(t *testing.T) {
	req := request.GetRatingBySourceTypeAndActorRequest{
		SourceType: "doctor",
		SourceUID:  "894",
	}
	resultRatings := []entity.RatingsCol{}
	filter := request.GetRatingBySourceTypeAndActorFilter{}
	publicRatingRepository.Mock.On("GetRatingsBySourceTypeAndActor", req.SourceType, req.SourceUID, filter).Return(resultRatings, nil)

	_, msg := svc.GetRatingBySourceTypeAndActor(req)
	assert.Equal(t, message.ErrNoData, msg)
}

func TestCreateRatingSubHelpfulSuccess(t *testing.T) {
	objectId, _ := primitive.ObjectIDFromHex(ratingSubId)
	objectRatingSubHelpfulId, _ := primitive.ObjectIDFromHex(ratingSubHelpId)
	counter := int64(4)

	input := request.CreateRatingSubHelpfulRequest{
		RatingSubmissionID: ratingSubId,
		UserID:             userId,
		UserIDLegacy:       userId,
		IPAddress:          ipaddress,
		UserAgent:          useragent,
	}

	ratingSubmission := entity.RatingSubmisson{
		ID:            objectId,
		RatingID:      ratingId,
		UserID:        &userId,
		UserIDLegacy:  &userId,
		Comment:       &comment,
		Value:         "85",
		IPAddress:     ipaddress,
		UserAgent:     useragent,
		SourceTransID: "",
		LikeCounter:   3,
	}

	ratingSubHelpful := entity.RatingSubHelpfulCol{
		ID:                 objectRatingSubHelpfulId,
		RatingSubmissionID: ratingSubId,
		UserID:             userId,
		UserIDLegacy:       userId,
		IPAddress:          ipaddress,
		UserAgent:          useragent,
	}

	ratingRepository.Mock.On("GetRatingSubmissionById", objectId).Return(ratingSubmission, nil).Once()
	publicRatingRepository.Mock.On("GetRatingSubHelpfulByRatingSubAndActor", input.RatingSubmissionID, input.UserIDLegacy).Return(nil, nil).Once()
	publicRatingRepository.Mock.On("CreateRatingSubHelpful", input).Return(ratingSubHelpful, nil).Once()
	publicRatingRepository.Mock.On("UpdateCounterRatingSubmission", objectId, counter).Return(nil).Once()

	_, msg := svc.CreateRatingSubHelpful(input)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestCreateRatingSubHelpfulRatingSubmissionNil(t *testing.T) {
	objectRatingSubmissionId, _ := primitive.ObjectIDFromHex(ratingSubIdFailed)

	input := request.CreateRatingSubHelpfulRequest{
		RatingSubmissionID: ratingSubIdFailed,
		UserID:             userId,
		UserIDLegacy:       userId,
		IPAddress:          ipaddress,
		UserAgent:          useragent,
	}

	ratingSubmission := entity.RatingSubmisson{
		ID:            objectRatingSubmissionId,
		RatingID:      ratingId,
		UserID:        &userId,
		UserIDLegacy:  &userId,
		Comment:       &comment,
		Value:         "85",
		IPAddress:     ipaddress,
		UserAgent:     useragent,
		SourceTransID: "",
		LikeCounter:   3,
	}

	ratingRepository.Mock.On("GetRatingSubmissionById", objectRatingSubmissionId).Return(ratingSubmission, nil).Once()

	_, msg := svc.CreateRatingSubHelpful(input)
	assert.Equal(t, message.FailedMsg, msg)
}

func TestCreateRatingSubHelpfulUpdateCounterFailed(t *testing.T) {
	objectId, _ := primitive.ObjectIDFromHex(ratingSubId)
	objectRatingSubHelpfulId, _ := primitive.ObjectIDFromHex(ratingSubHelpId)

	input := request.CreateRatingSubHelpfulRequest{
		RatingSubmissionID: ratingSubId,
		UserID:             userId,
		UserIDLegacy:       userId,
		IPAddress:          ipaddress,
		UserAgent:          useragent,
	}

	ratingSubmission := entity.RatingSubmisson{
		ID:            objectId,
		RatingID:      ratingid,
		UserID:        &userId,
		UserIDLegacy:  &userId,
		Comment:       &comment,
		Value:         "85",
		IPAddress:     ipaddress,
		UserAgent:     useragent,
		SourceTransID: "",
		LikeCounter:   3,
	}

	ratingSubHelpful := entity.RatingSubHelpfulCol{
		ID:                 objectRatingSubHelpfulId,
		RatingSubmissionID: ratingSubId,
		UserID:             userId,
		UserIDLegacy:       userId,
		IPAddress:          ipaddress,
		UserAgent:          useragent,
	}

	ratingRepository.Mock.On("GetRatingSubmissionById", objectId).Return(ratingSubmission, nil).Once()
	publicRatingRepository.Mock.On("GetRatingSubHelpfulByRatingSubAndActor", input.RatingSubmissionID, input.UserIDLegacy).Return(nil, nil).Once()
	publicRatingRepository.Mock.On("CreateRatingSubHelpful", input).Return(ratingSubHelpful, nil).Once()
	publicRatingRepository.Mock.On("UpdateCounterRatingSubmission", objectId, ratingSubmission.LikeCounter).Return(errors.New("error")).Once()

	_, msg := svc.CreateRatingSubHelpful(input)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetRatingSummaryBySourceType(t *testing.T) {
	idObj, _ := primitive.ObjectIDFromHex(idDummy1)
	ratingTypeObj, _ := primitive.ObjectIDFromHex(ratingid)

	ratingDatas := []entity.RatingsCol{
		{
			ID:           idObj,
			Name:         "Rating 1 Doctor A",
			Description:  &description,
			SourceUid:    "3310",
			SourceType:   requestSummary.SourceType,
			RatingType:   ratingType,
			RatingTypeId: ratingid,
		},
	}
	ratingSubDatas := []entity.RatingSubmisson{
		{
			ID:            idObj,
			RatingID:      idDummy1,
			UserID:        &userId,
			UserIDLegacy:  &userId,
			Comment:       &comment,
			Value:         "90",
			IPAddress:     ipaddress,
			UserAgent:     useragent,
			SourceTransID: "",
			LikeCounter:   5,
		},
	}
	ratingFormula := entity.RatingFormulaCol{
		ID:           idObj,
		SourceType:   "doctor",
		Formula:      "(9000 + total_rating_point) / (100 + total_user_count)",
		RatingTypeId: ratingid,
		RatingType:   ratingType,
	}
	paginationResult := base.Pagination{
		Records:      1,
		Limit:        10,
		Page:         1,
		TotalRecords: 1,
	}
	publicRatingRepository.Mock.On("GetPublicRatingsByParams", requestSummary.Limit, requestSummary.Page, "updated_at", filterSummary).Return(ratingDatas, &paginationResult, nil).Once()
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", ratingTypeObj).Return(nil, mongo.ErrNoDocuments).Once()
	publicRatingRepository.Mock.On("GetRatingSubsByRatingId", idObj.Hex()).Return(ratingSubDatas, nil).Once()
	publicRatingRepository.Mock.On("GetRatingFormulaByRatingTypeIdAndSourceType", ratingid, requestSummary.SourceType).Return(&ratingFormula, nil).Once()

	result, pagination, msg := publicRatingService.GetListRatingSummaryBySourceType(requestSummary)
	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 1000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be success")
	assert.Equal(t, 1, len(result), "Count of list kd must be 1")
	assert.Equal(t, int64(1), pagination.Records, "Total record must be 1")
}

func TestGetRatingSummaryBySourceTypeErrEmptyRating(t *testing.T) {
	ratingDatas := []entity.RatingsCol{}
	paginationResult := base.Pagination{
		Records:      0,
		Limit:        10,
		Page:         1,
		TotalRecords: 0,
	}
	publicRatingRepository.Mock.On("GetPublicRatingsByParams", requestSummary.Limit, requestSummary.Page, "updated_at", filterSummary).Return(ratingDatas, &paginationResult, errors.New("error")).Once()

	result, pagination, msg := publicRatingService.GetListRatingSummaryBySourceType(requestSummary)
	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 1000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be success")
	assert.Equal(t, 0, len(result), "Count of list kd must be 0")
	assert.Equal(t, int64(0), pagination.Records, "Total record must be 0")
}

func TestGetRatingSummaryBySourceTypeFailedGetRating(t *testing.T) {
	var request = request.GetPublicListRatingSummaryRequest{
		SourceType: "doctor",
		Sort:       "failed",
		Dir:        "desc",
		Page:       1,
		Limit:      50,
	}
	ratingDatas := []entity.RatingsCol{}
	paginationResult := base.Pagination{
		Records:      0,
		Limit:        10,
		Page:         1,
		TotalRecords: 1,
	}
	publicRatingRepository.Mock.On("GetPublicRatingsByParams", requestSummary.Limit, requestSummary.Page, "Collection tidak ditemukan", filterSummary).Return(ratingDatas, &paginationResult, errors.New("error")).Once()

	_, _, msg := publicRatingService.GetListRatingSummaryBySourceType(request)
	assert.Equal(t, message.RecordNotFound.Code, msg.Code, "Code must be 412002")
	assert.Equal(t, message.RecordNotFound.Message, msg.Message, "Message must be failed")
}

func TestGetRatingSummaryBySourceTypeErrGetRatingSubmission(t *testing.T) {
	idDummy1, _ := primitive.ObjectIDFromHex(failedId)
	ratingTypeObj, _ := primitive.ObjectIDFromHex(ratingid)
	ratingDatas := []entity.RatingsCol{
		{
			ID:           idDummy1,
			Name:         "Rating 1 Doctor A",
			Description:  &description,
			SourceUid:    "3310",
			SourceType:   requestSummary.SourceType,
			RatingType:   ratingType,
			RatingTypeId: ratingid,
		},
	}
	paginationResult := base.Pagination{
		Records:      1,
		Limit:        10,
		Page:         1,
		TotalRecords: 1,
	}
	publicRatingRepository.Mock.On("GetPublicRatingsByParams", requestSummary.Limit, requestSummary.Page, "updated_at", filterSummary).Return(ratingDatas, &paginationResult, nil).Once()
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", ratingTypeObj).Return(nil, mongo.ErrNoDocuments).Once()
	publicRatingRepository.Mock.On("GetRatingSubsByRatingId", failedId).Return(nil, errors.New("error")).Once()

	_, _, msg := publicRatingService.GetListRatingSummaryBySourceType(requestSummary)
	assert.Equal(t, message.ErrFailedSummaryRatingNumeric.Code, msg.Code, "Code must be 412002")
	assert.Equal(t, message.ErrFailedSummaryRatingNumeric.Message, msg.Message, "Message must be Failed to summary rating numeric")
}

func TestGetRatingSummaryBySourceTypeErrFailedCalculate(t *testing.T) {
	failID := "62c3e57b457ed515928c3690"
	idObj, _ := primitive.ObjectIDFromHex(idDummy1)
	ratingTypeObj, _ := primitive.ObjectIDFromHex(failID)

	ratingDatas := []entity.RatingsCol{
		{
			ID:           idObj,
			Name:         "Rating 1 Doctor A",
			Description:  &description,
			SourceUid:    "3310",
			SourceType:   requestSummary.SourceType,
			RatingType:   ratingType,
			RatingTypeId: failID,
		},
	}
	ratingSubDatas := []entity.RatingSubmisson{
		{
			ID:       idObj,
			RatingID: idDummy1,
			Value:    "k",
		},
	}
	paginationResult := base.Pagination{
		Records:      1,
		Limit:        10,
		Page:         1,
		TotalRecords: 1,
	}
	publicRatingRepository.Mock.On("GetPublicRatingsByParams", requestSummary.Limit, requestSummary.Page, "updated_at", filterSummary).Return(ratingDatas, &paginationResult, errors.New("error")).Once()
	ratingRepository.Mock.On("GetRatingTypeLikertByIdAndStatus", ratingTypeObj).Return(nil, mongo.ErrNoDocuments).Once()
	publicRatingRepository.Mock.On("GetRatingSubsByRatingId", idObj.Hex()).Return(ratingSubDatas, nil).Once()
	publicRatingRepository.Mock.On("GetRatingFormulaByRatingTypeIdAndSourceType", failID, requestSummary.SourceType).Return(nil, nil).Once()

	_, _, msg := publicRatingService.GetListRatingSummaryBySourceType(requestSummary)
	assert.Equal(t, message.ErrFailedSummaryRatingNumeric.Code, msg.Code, "Code must be 412002")
	assert.Equal(t, message.ErrFailedSummaryRatingNumeric.Message, msg.Message, "Message must be Failed to summary rating numeric")
}

func TestGetRatingSubmissionBySourceTypeAndUID(t *testing.T) {
	idDummy1, _ := primitive.ObjectIDFromHex(idDummy1)
	idDummy2, _ := primitive.ObjectIDFromHex(idDummy2)

	var filterSubmission = request.FilterRatingSubmission{
		RatingID: []string{idDummy1.Hex(), idDummy2.Hex()},
	}

	ratingDatas := []entity.RatingsCol{
		{
			ID:           idDummy1,
			Name:         "Rating 1 Doctor A",
			Description:  &description,
			SourceUid:    "3310",
			SourceType:   requestSummary.SourceType,
			RatingType:   ratingType,
			RatingTypeId: ratingid,
		},
		{
			ID:           idDummy2,
			Name:         "Rating 1 Doctor B",
			Description:  &description,
			SourceUid:    "3311",
			SourceType:   requestSummary.SourceType,
			RatingType:   ratingType,
			RatingTypeId: ratingid,
		},
	}
	ratingSubDatas := []entity.RatingSubmisson{
		{
			ID:            idDummy1,
			RatingID:      idDummy1.Hex(),
			UserID:        &userId,
			UserIDLegacy:  &userId,
			DisplayName:   &displayName,
			Comment:       &comment,
			SourceTransID: sourceTransId,
			LikeCounter:   5,
			IsAnonymous:   anonym,
		},
	}
	paginationResult := base.Pagination{
		Records:      1,
		Limit:        10,
		Page:         1,
		TotalRecords: 1,
	}
	publicRatingRepository.Mock.On("GetListRatingBySourceTypeAndUID", requestSubmission.SourceType, requestSubmission.SourceUID).Return(ratingDatas, nil).Once()
	publicRatingRepository.Mock.On("GetPublicRatingSubmissions", requestSubmission.Limit, requestSubmission.Page, "created_at", filterSubmission).Return(ratingSubDatas, &paginationResult, nil).Once()
	ratingRepository.Mock.On("GetRatingById", idDummy1).Return(&ratingDatas[0], nil)

	result, pagination, msg := publicRatingService.GetListRatingSubmissionBySourceTypeAndUID(requestSubmission)
	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 1000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be success")
	assert.Equal(t, 1, len(result), "Count of list kd must be 1")
	assert.Equal(t, int64(1), pagination.Records, "Total record must be 1")
}

func TestGetRatingSubmissionBySourceTypeAndUIDEmptyList(t *testing.T) {
	idDummy1, _ := primitive.ObjectIDFromHex(idDummy1)
	idDummy2, _ := primitive.ObjectIDFromHex(idDummy2)

	var filterSubmission = request.FilterRatingSubmission{
		RatingID: []string{idDummy1.Hex(), idDummy2.Hex()},
	}

	ratingDatas := []entity.RatingsCol{
		{
			ID:           idDummy1,
			Name:         "Rating 1 Doctor A",
			Description:  &description,
			SourceUid:    "3310",
			SourceType:   requestSummary.SourceType,
			RatingType:   ratingType,
			RatingTypeId: ratingid,
		},
		{
			ID:           idDummy2,
			Name:         "Rating 1 Doctor B",
			Description:  &description,
			SourceUid:    "3311",
			SourceType:   requestSummary.SourceType,
			RatingType:   ratingType,
			RatingTypeId: ratingid,
		},
	}
	ratingSubDatas := []entity.RatingSubmisson{}
	paginationResult := base.Pagination{
		Records:      0,
		Limit:        10,
		Page:         1,
		TotalRecords: 0,
	}
	publicRatingRepository.Mock.On("GetListRatingBySourceTypeAndUID", requestSubmission.SourceType, requestSubmission.SourceUID).Return(ratingDatas, nil).Once()
	publicRatingRepository.Mock.On("GetPublicRatingSubmissions", requestSubmission.Limit, requestSubmission.Page, "created_at", filterSubmission).Return(ratingSubDatas, &paginationResult, nil).Once()

	result, pagination, msg := publicRatingService.GetListRatingSubmissionBySourceTypeAndUID(requestSubmission)
	assert.Equal(t, message.ErrNoData.Code, msg.Code, "Code must be 212004")
	assert.Equal(t, message.ErrNoData.Message, msg.Message, "Message must be erro data not found")
	assert.Equal(t, 0, len(result), "Count of list kd must be 0")
	assert.Equal(t, int64(0), pagination.Records, "Total record must be 0")
}

func TestGetRatingSubmissionBySourceTypeAndUIDErrGetRating(t *testing.T) {
	message := "Cannot find rating with params SourceType :" + requestSubmission.SourceType + ", SourceUid:" + requestSubmission.SourceUID
	publicRatingRepository.Mock.On("GetListRatingBySourceTypeAndUID", requestSubmission.SourceType, requestSubmission.SourceUID).Return([]entity.RatingsCol{}, errors.New("error")).Once()

	_, _, msg := publicRatingService.GetListRatingSubmissionBySourceTypeAndUID(requestSubmission)
	assert.Equal(t, 412002, msg.Code, "Code must be 412002")
	assert.Equal(t, message, msg.Message, "Message must be "+message)
}
