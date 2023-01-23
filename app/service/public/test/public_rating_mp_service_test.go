package publictest

import (
	"context"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	publicrequest "go-klikdokter/app/model/request/public"
	publicresponse "go-klikdokter/app/model/response/public"
	"go-klikdokter/app/repository/public/public_repository_mock"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service/public"
	"go-klikdokter/helper/message"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ratingMpRepository = &repository_mock.RatingMpRepository{Mock: mock.Mock{}}
var publicRatingMpRepository = &public_repository_mock.PublicRatingMpRepository{Mock: mock.Mock{}}
var publicRatingMpService = publicservice.NewPublicRatingMpService(logger, ratingMpRepository, publicRatingMpRepository)

func init() {
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}
}

var (
	sourceTransIdMp = "ORDER-1234"
)
var idDummy1Obj, _ = primitive.ObjectIDFromHex(idDummy1)

var requestSummaryMp = publicrequest.GetPublicListRatingSummaryRequest{
	SourceType: "product",
	Filter:     `{"source_uid": ["1234"], "rating_type": ["rating_for_product"]}`,
	Sort:       "",
	Dir:        "desc",
	Page:       1,
	Limit:      50,
}

var requestSubmissionMp = publicrequest.GetPublicListRatingSubmissionRequest{
	SourceType: "product",
	SourceUID:  "1234",
	Sort:       "created_at",
	Dir:        "desc",
	Page:       1,
	Limit:      50,
}

var filterSummaryMp = publicrequest.FilterRatingSummary{
	SourceType: requestSummaryMp.SourceType,
	SourceUid:  []string{"1234"},
	RatingType: []string{"rating_for_product"},
}

var requestSubmissionMpByID = publicrequest.GetPublicListRatingSubmissionByIDRequest{
	Filter: `{"rating_subs_id": ["1234"]}`,
	Source: "mp",
	Limit:  50,
	Page:   1,
	Sort:   "created_at",
	Dir:    "1",
}

var requestSummaryMpDetail = publicrequest.PublicGetListDetailRatingSummaryRequest{
	Filter:     `{"source_uid": ["1234"]}`,
	SourceType: "product",
	Limit:      50,
	Page:       1,
	Sort:       "created_at",
	Dir:        "1",
}

var requestSummaryStoreProduct = publicrequest.PublicGetRatingSummaryStoreProductRequest{
	Filter: `{"store_uid": ["1"]}`,
	Limit:  50,
	Page:   1,
	Sort:   "created_at",
	Dir:    "1",
}

var ratingSubMpDatas = []entity.RatingSubmissionMp{
	{
		ID:            idDummy1Obj,
		UserID:        &userId,
		UserIDLegacy:  &userId,
		Comment:       &comment,
		Value:         "4",
		IPAddress:     ipaddress,
		UserAgent:     useragent,
		SourceTransID: sourceTransIdMp,
	},
}

func TestGetRatingSummaryMpBySourceType(t *testing.T) {
	idObj, _ := primitive.ObjectIDFromHex(idDummy1)
	ratingSubDatas := []entity.RatingSubmissionMp{
		{
			ID:            idObj,
			UserID:        &userId,
			UserIDLegacy:  &userId,
			Comment:       &comment,
			Value:         "4",
			IPAddress:     ipaddress,
			UserAgent:     useragent,
			SourceTransID: sourceTransIdMp,
		},
	}
	ratingFormulaMp := entity.RatingFormulaCol{
		ID:           idObj,
		SourceType:   "product",
		Formula:      "(sum / count) / 1",
		RatingTypeId: ratingid,
		RatingType:   ratingType,
	}
	paginationResult := base.Pagination{
		Records:      1,
		Limit:        10,
		Page:         1,
		TotalRecords: 1,
	}
	ratingSubmissionGroupBySource := publicresponse.PublicRatingSubGroupBySourceMp{
		ID: struct {
			SourceUID  string `json:"source_uid" bson:"source_uid"`
			SourceType string `json:"source_type" bson:"source_type"`
		}(struct {
			SourceUID  string
			SourceType string
		}{SourceUID: "1234", SourceType: "product"}),
		RatingSubmissionsMp: ratingSubDatas,
	}
	sumCountRatingSummary := publicresponse.PublicSumCountRatingSummaryMp{
		Sum:      10,
		Count:    3,
		Comments: []string{},
	}

	publicRatingMpRepository.Mock.On("GetPublicRatingSubmissionsGroupBySource", requestSummaryMp.Limit, requestSummaryMp.Page, -1, "updated_at", filterSummaryMp).
		Return([]publicresponse.PublicRatingSubGroupBySourceMp{ratingSubmissionGroupBySource}, &paginationResult, nil).Once()
	publicRatingMpRepository.Mock.On("GetSumCountRatingSubsBySource", ratingSubmissionGroupBySource.ID.SourceUID, ratingSubmissionGroupBySource.ID.SourceType).Return(&sumCountRatingSummary, nil).Once()
	publicRatingMpRepository.Mock.On("GetRatingFormulaBySourceType", requestSummaryMp.SourceType).Return(&ratingFormulaMp, nil).Once()

	result, pagination, msg := publicRatingMpService.GetListRatingSummaryBySourceType(requestSummaryMp)
	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 1000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be success")
	assert.Equal(t, 1, len(result), "Count of list kd must be 1")
	assert.Equal(t, int64(1), pagination.Records, "Total record must be 1")
}

func TestGetRatingSubmissionMpBySourceTypeAndUID(t *testing.T) {
	idDummy1, _ := primitive.ObjectIDFromHex(idDummy1)

	var filterSubmission = publicrequest.FilterRatingSubmissionMp{
		SourceUID:  "1234",
		SourceType: "product",
	}

	ratingSubDatas := []entity.RatingSubmissionMp{
		{
			ID:            idDummy1,
			UserID:        &userId,
			UserIDLegacy:  &userId,
			DisplayName:   &displayName,
			Comment:       &comment,
			SourceTransID: sourceTransIdMp,
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
	publicRatingMpRepository.Mock.On("GetPublicRatingSubmissions", requestSubmission.Limit, requestSubmission.Page, -1, "created_at", filterSubmission).
		Return(ratingSubDatas, &paginationResult, nil).Once()

	result, pagination, msg := publicRatingMpService.GetListRatingSubmissionBySourceTypeAndUID(requestSubmissionMp)
	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 1000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be success")
	assert.Equal(t, 1, len(result), "Count of list kd must be 1")
	assert.Equal(t, int64(1), pagination.Records, "Total record must be 1")
}

func TestGetListRatingSubmissionMpByID(t *testing.T) {
	idObj, _ := primitive.ObjectIDFromHex(idDummy1)
	ratingSubDatas := []entity.RatingSubmissionMp{
		{
			ID:            idObj,
			UserID:        &userId,
			UserIDLegacy:  &userId,
			Comment:       &comment,
			Value:         "4",
			IPAddress:     ipaddress,
			UserAgent:     useragent,
			SourceTransID: sourceTransIdMp,
			Media: []entity.MediaObj{
				{
					UID:       "1111",
					MediaPath: "aaaaa",
				},
			},
		},
	}
	paginationResult := base.Pagination{
		Records:      1,
		Limit:        10,
		Page:         1,
		TotalRecords: 1,
	}
	filterRatingSubs := publicrequest.FilterRatingSubmissionMp{
		RatingSubsID: []string{"1234"},
	}

	publicRatingMpRepository.Mock.On("GetPublicRatingSubmissionsCustom", requestSubmissionMpByID.Limit, requestSubmissionMpByID.Page, -1, "created_at", filterRatingSubs, requestSubmissionMpByID.Source).
		Return(ratingSubDatas, &paginationResult, nil).Once()
	ctx := context.Background()
	result, pagination, msg, _ := publicRatingMpService.GetListRatingSubmissionByID(ctx, requestSubmissionMpByID)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 1000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be success")
	assert.Equal(t, 1, len(result), "Count of list kd must be 1")
	assert.Equal(t, int64(1), pagination.Records, "Total record must be 1")
}

func TestGetListDetailRatingSummaryMpBySourceType(t *testing.T) {
	idObj, _ := primitive.ObjectIDFromHex(idDummy1)
	paginationResult := base.Pagination{
		Records:      1,
		Limit:        10,
		Page:         1,
		TotalRecords: 1,
	}
	ratingSubmissionGroupBySource := publicresponse.PublicRatingSubGroupBySourceMp{
		ID: struct {
			SourceUID  string `json:"source_uid" bson:"source_uid"`
			SourceType string `json:"source_type" bson:"source_type"`
		}(struct {
			SourceUID  string
			SourceType string
		}{SourceUID: "1234", SourceType: "product"}),
		RatingSubmissionsMp: ratingSubMpDatas,
	}

	var filterListDetailRatingSummaryMp = publicrequest.FilterRatingSummary{
		SourceType: requestSummaryMp.SourceType,
		SourceUid:  []string{"1234"},
	}

	ratingFormulaMp := entity.RatingFormulaCol{
		ID:           idObj,
		SourceType:   "product",
		Formula:      "(sum / count) / 1",
		RatingTypeId: ratingid,
		RatingType:   ratingType,
	}

	publicRatingMpRepository.Mock.On("GetPublicRatingSubmissionsGroupBySource", requestSummaryMp.Limit, requestSummaryMp.Page, -1, "created_at", filterListDetailRatingSummaryMp).
		Return([]publicresponse.PublicRatingSubGroupBySourceMp{ratingSubmissionGroupBySource}, &paginationResult, nil).Once()
	publicRatingMpRepository.Mock.On("GetRatingFormulaBySourceType", requestSummaryMp.SourceType).Return(&ratingFormulaMp, nil).Once()

	result, pagination, msg := publicRatingMpService.GetListDetailRatingSummaryBySourceType(requestSummaryMpDetail)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 1000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be success")
	assert.Equal(t, 1, len(result), "Count of list kd must be 1")
	assert.Equal(t, int64(1), pagination.Records, "Total record must be 1")
}

func TestGetRatingSummaryStoreProduct(t *testing.T) {
	idObj, _ := primitive.ObjectIDFromHex(idDummy1)
	paginationResult := base.Pagination{
		Records:      1,
		Limit:        10,
		Page:         1,
		TotalRecords: 1,
	}

	ratingSubmissionGroupByStoreSource := publicresponse.PublicRatingSubGroupByStoreSourceMp{
		ID: struct {
			StoreUID   string `json:"store_uid" bson:"store_uid"`
			SourceType string `json:"source_type" bson:"source_type"`
		}(struct {
			StoreUID   string
			SourceType string
		}{StoreUID: "1", SourceType: "product"}),
		RatingSubmissionsMp: ratingSubMpDatas,
	}

	var filter = publicrequest.FilterRatingSummary{
		SourceType: "product",
		StoreUID:   []string{"1"},
	}

	ratingFormulaMp := entity.RatingFormulaCol{
		ID:           idObj,
		SourceType:   "product",
		Formula:      "(sum / count) / 1",
		RatingTypeId: ratingid,
		RatingType:   ratingType,
	}

	publicRatingMpRepository.Mock.On("GetPublicRatingSubmissionsGroupByStoreSource", requestSummaryStoreProduct.Limit, requestSummaryStoreProduct.Page, -1, "created_at", filter).
		Return([]publicresponse.PublicRatingSubGroupByStoreSourceMp{ratingSubmissionGroupByStoreSource}, &paginationResult, nil).Once()
	publicRatingMpRepository.Mock.On("GetRatingFormulaBySourceType", "product").Return(&ratingFormulaMp, nil).Once()

	result, pagination, msg := publicRatingMpService.GetRatingSummaryStoreProduct(context.TODO(), requestSummaryStoreProduct)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 1000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be success")
	assert.Equal(t, 1, len(result), "Count of list kd must be 1")
	assert.Equal(t, int64(1), pagination.Records, "Total record must be 1")
}
