package publicservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-klikdokter/app/model/base"
	publicrequest "go-klikdokter/app/model/request/public"
	"go-klikdokter/app/model/response"
	publicresponse "go-klikdokter/app/model/response/public"
	"go-klikdokter/app/repository"
	publicrepository "go-klikdokter/app/repository/public"
	"go-klikdokter/helper/config"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"
	"go-klikdokter/helper/thumbor"
	"go-klikdokter/pkg/util"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
	"github.com/vjeantet/govaluate"
	"go.mongodb.org/mongo-driver/mongo"
)

type PublicRatingMpService interface {
	GetListRatingSubmissionBySourceTypeAndUID(input publicrequest.GetPublicListRatingSubmissionRequest) ([]publicresponse.PublicRatingSubmissionMpResponse, *base.Pagination, message.Message)
	GetListRatingSummaryBySourceType(input publicrequest.GetPublicListRatingSummaryRequest) ([]publicresponse.PublicRatingSummaryMpResponse, message.Message)
	GetListRatingSubmissionByID(ctx context.Context, input publicrequest.GetPublicListRatingSubmissionByIDRequest) ([]publicresponse.PublicRatingSubmissionMpResponse, *base.Pagination, message.Message, interface{})
	GetListDetailRatingSummaryBySourceType(input publicrequest.PublicGetListDetailRatingSummaryRequest) ([]publicresponse.PublicRatingSummaryListDetailResponse, message.Message)
	GetRatingSummaryStoreProduct(ctx context.Context, input publicrequest.PublicGetRatingSummaryStoreProductRequest) ([]publicresponse.RatingSummaryStoreProductNumeric, message.Message)
}

type publicRatingMpServiceImpl struct {
	logger             log.Logger
	ratingMpRepo       repository.RatingMpRepository
	publicRatingMpRepo publicrepository.PublicRatingMpRepository
}

func NewPublicRatingMpService(
	lg log.Logger,
	rrmp repository.RatingMpRepository,
	prr publicrepository.PublicRatingMpRepository,
) PublicRatingMpService {
	return &publicRatingMpServiceImpl{lg, rrmp, prr}
}

func (s *publicRatingMpServiceImpl) GetListRatingSubmissionBySourceTypeAndUID(input publicrequest.GetPublicListRatingSubmissionRequest) ([]publicresponse.PublicRatingSubmissionMpResponse, *base.Pagination, message.Message) {
	results := []publicresponse.PublicRatingSubmissionMpResponse{}
	timezone := config.GetConfigString(viper.GetString("util.timezone"))
	avatarDefault := config.GetConfigString(viper.GetString("image.default-avatar"))
	var dir int
	if input.Dir == "asc" {
		dir = 1
	} else {
		dir = -1
	}
	// Set default value
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 50
	}
	if input.Sort == "" {
		input.Sort = "created_at"
	}

	// Unmarshal filter params
	filterRatingSubs := publicrequest.FilterRatingSubmissionMp{}
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filterRatingSubs)
		if errMarshal != nil {
			return nil, nil, message.ErrUnmarshalFilterListRatingRequest
		}
	}
	// new filter by source_uid and source_type
	filterRatingSubs.SourceUID = input.SourceUID
	filterRatingSubs.SourceType = input.SourceType

	// validation filter date
	if err := filterRatingSubs.ValidateFormatDate(); err != nil {
		return nil, nil, message.ErrInvalidDate
	}
	if filterRatingSubs.StartDate != "" && filterRatingSubs.EndDate != "" {
		startDate, _ := util.ConvertToDateTime(filterRatingSubs.StartDate + " 00:00:00")
		endDate, _ := util.ConvertToDateTime(filterRatingSubs.EndDate + " 23:59:59")
		if errD := endDate.Before(startDate); errD == true {
			return nil, nil, message.ErrRangeDate
		}
	}

	// Get Rating Submission
	ratingSubs, pagination, err := s.publicRatingMpRepo.GetPublicRatingSubmissions(input.Limit, input.Page, dir, input.Sort, filterRatingSubs)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return results, pagination, message.ErrNoData
		}
		return nil, nil, message.FailedMsg
	}
	if len(ratingSubs) <= 0 {
		return results, pagination, message.ErrNoData
	}

	for _, v := range ratingSubs {
		if v.Avatar == "" {
			v.Avatar = avatarDefault
		}
		Loc, _ := time.LoadLocation(timezone)

		// Masking Anonym Display Name
		displayName := ""
		if v.IsAnonymous {
			displayName = util.MaskDisplayName(*v.DisplayName)
		} else {
			displayName = *v.DisplayName
		}
		// update media_path from null to empty array
		// create thumbor response
		var mediaResponse = []response.MediaObjResponse{}
		for _, value := range v.Media {
			mediaObjResponse := response.MediaObjResponse{
				UID:        value.UID,
				MediaPath:  value.MediaPath,
				MediaImage: thumbor.GetNewThumborImagesOriginal(value.MediaPath),
			}
			mediaResponse = append(mediaResponse, mediaObjResponse)
		}

		results = append(results, publicresponse.PublicRatingSubmissionMpResponse{
			ID:            v.ID,
			SourceType:    v.SourceType,
			SourceUID:     v.SourceUID,
			StoreUID:      v.StoreUID,
			UserID:        v.UserID,
			UserIDLegacy:  v.UserIDLegacy,
			DisplayName:   displayName,
			Avatar:        v.Avatar,
			Comment:       v.Comment,
			SourceTransID: v.SourceTransID,
			LikeCounter:   v.LikeCounter,
			Value:         strconv.Itoa(v.Value),
			LikeByMe:      false,
			Media:         mediaResponse,
			IsWithMedia:   v.IsWithMedia,
			CreatedAt:     v.CreatedAt.In(Loc),
			Reply:         v.Reply,
			ReplyBy:       v.ReplyBy,
		})

	}
	return results, pagination, message.SuccessMsg
}

func (s *publicRatingMpServiceImpl) GetListRatingSummaryBySourceType(input publicrequest.GetPublicListRatingSummaryRequest) ([]publicresponse.PublicRatingSummaryMpResponse, message.Message) {
	results := []publicresponse.PublicRatingSummaryMpResponse{}
	input.MakeDefaultValueIfEmpty()
	filter := publicrequest.FilterRatingSummary{}
	filter.SourceType = input.SourceType
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filter)
		if errMarshal != nil {
			return nil, message.ErrUnmarshalFilterListRatingRequest
		}
	}
	if errValidate := filter.ValidateSourceUID(); errValidate != nil {
		return nil, *errValidate
	}

	ratingSub, err := s.publicRatingMpRepo.GetPublicRatingSubmissionsGroupBySource(filter)
	if err != nil {
		return nil, message.RecordNotFound
	}
	if len(ratingSub) <= 0 {
		return results, message.SuccessMsg
	}

	for _, args := range ratingSub {
		data, err := s.summaryRatingNumeric(args.ID.SourceUID, input.SourceType)
		if err != nil {
			return nil, message.ErrFailedSummaryRatingNumeric
		}
		results = append(results, *data)
	}
	return results, message.SuccessMsg
}

func (s *publicRatingMpServiceImpl) summaryRatingNumeric(sourceUID string, sourceType string) (*publicresponse.PublicRatingSummaryMpResponse, error) {
	sumCountRatingSubs, err := s.publicRatingMpRepo.GetSumCountRatingSubsBySource(sourceUID, sourceType)
	if err != nil {
		return nil, err
	}

	if sumCountRatingSubs == nil {
		return nil, errors.New("data RatingSubmission not found")
	}

	formulaRating, err := s.publicRatingMpRepo.GetRatingFormulaBySourceType(sourceType)
	if err != nil {
		return nil, err
	}
	if formulaRating == nil {
		return nil, errors.New("Formula rating not found, for source_type: " + sourceType)
	}

	if formulaRating.Formula != "" {
		ratingSummary, err := calculateRatingMpValue(sourceUID, formulaRating.Formula, sumCountRatingSubs)
		if err != nil {
			return nil, err
		}
		ratingSummary.MaximumValue = global.GetMaximumValueBySourceType(sourceType)

		data := publicresponse.PublicRatingSummaryMpResponse{
			SourceUid:     sourceUID,
			SourceType:    sourceType,
			RatingSummary: ratingSummary,
		}
		return &data, nil
	} else {
		return nil, errors.New("formula string is empty")
	}
}

// swagger:route GET /public/rating-submissions-by-id PublicRating GetPublicListRatingSubmissionByIDRequest
// Get Rating Submission By ID
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *publicRatingMpServiceImpl) GetListRatingSubmissionByID(ctx context.Context, input publicrequest.GetPublicListRatingSubmissionByIDRequest) ([]publicresponse.PublicRatingSubmissionMpResponse, *base.Pagination, message.Message, interface{}) {
	err := input.ValidateFilterAndSource()
	if err != nil {
		return nil, nil, message.ErrReqParam, err
	}
	errMsg := make(map[string]interface{})
	avatarDefault := config.GetConfigString(viper.GetString("image.default-avatar"))
	timezone := config.GetConfigString(viper.GetString("util.timezone"))
	var result []publicresponse.PublicRatingSubmissionMpResponse
	// Unmarshal filter params
	filterRatingSubs := publicrequest.FilterRatingSubmissionMp{}
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filterRatingSubs)
		if errMarshal != nil {
			return result, nil, message.ErrUnmarshalFilterListRatingRequest, nil
		}
	}
	// Get Rating Submission
	var dir int
	if input.Dir == "asc" {
		dir = 1
	} else {
		dir = -1
	}
	// Set default value
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 100
	}
	if input.Sort == "" {
		input.Sort = "created_at"
	}
	ratingSubs, pagination, err := s.publicRatingMpRepo.GetPublicRatingSubmissionsCustom(input.Limit, input.Page, dir, input.Sort, filterRatingSubs, input.Source)
	if err != nil {
		errMsg["rating"] = "Error Not Found"
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, pagination, message.ErrNoData, errMsg
		}
		return result, pagination, message.FailedMsg, errMsg
	}
	if len(ratingSubs) <= 0 {
		return result, pagination, message.SuccessMsg, nil
	}

	for _, v := range ratingSubs {
		if v.Avatar == "" {
			v.Avatar = avatarDefault
		}
		Loc, _ := time.LoadLocation(timezone)

		// Masking Anonym Display Name
		displayName := ""
		if v.DisplayName != nil {
			if v.IsAnonymous {
				displayName = util.MaskDisplayName(*v.DisplayName)
			} else {
				displayName = *v.DisplayName
			}
		}

		// create thumbor response
		var mediaResponse = []response.MediaObjResponse{}
		for _, value := range v.Media {
			mediaObjResponse := response.MediaObjResponse{
				UID:        value.UID,
				MediaPath:  value.MediaPath,
				MediaImage: thumbor.GetNewThumborImagesOriginal(value.MediaPath),
			}
			mediaResponse = append(mediaResponse, mediaObjResponse)
		}
		result = append(result, publicresponse.PublicRatingSubmissionMpResponse{
			ID:            v.ID,
			SourceType:    v.SourceType,
			SourceUID:     v.SourceUID,
			StoreUID:      v.StoreUID,
			UserID:        v.UserID,
			UserIDLegacy:  v.UserIDLegacy,
			DisplayName:   displayName,
			Avatar:        v.Avatar,
			Comment:       v.Comment,
			SourceTransID: v.SourceTransID,
			LikeCounter:   v.LikeCounter,
			Value:         strconv.Itoa(v.Value),
			LikeByMe:      false,
			IsWithMedia:   v.IsWithMedia,
			Media:         mediaResponse,
			CreatedAt:     v.CreatedAt.In(Loc),
			Reply:         v.Reply,
			ReplyBy:       v.ReplyBy,
		})

	}
	return result, pagination, message.SuccessMsg, nil
}

func calculateRatingMpValue(sourceUID, formula string, sumCountRatingSubs *publicresponse.PublicSumCountRatingSummaryMp) (publicresponse.RatingSummaryMpNumeric, error) {
	result := publicresponse.RatingSummaryMpNumeric{}
	result.SourceUID = sourceUID
	result.TotalReviewer = sumCountRatingSubs.Count
	result.TotalValue = "0"

	if len(sumCountRatingSubs.Comments) > 0 {
		for _, c := range sumCountRatingSubs.Comments {
			if strings.TrimSpace(c) != "" {
				result.TotalComment++
			}
		}
	} else {
		result.TotalComment = result.TotalReviewer
	}

	if formula != "" {
		expression, err := govaluate.NewEvaluableExpression(formula)
		if err != nil {
			return result, err
		}
		parameters := make(map[string]interface{}, 8)
		parameters["sum"] = sumCountRatingSubs.Sum
		parameters["count"] = sumCountRatingSubs.Count

		finalCalc, err := expression.Evaluate(parameters)
		if err != nil {
			return result, err
		}
		result.TotalValue = fmt.Sprintf("%.1f", finalCalc)
	}

	return result, nil
}

// swagger:route GET /public/ratings-summary/detail/{source_type} PublicRating PublicGetListDetailRatingSummaryRequest
// Get List Detail Rating Summary
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *publicRatingMpServiceImpl) GetListDetailRatingSummaryBySourceType(input publicrequest.PublicGetListDetailRatingSummaryRequest) ([]publicresponse.PublicRatingSummaryListDetailResponse, message.Message) {
	results := []publicresponse.PublicRatingSummaryListDetailResponse{}
	input.MakeDefaultValueIfEmpty()

	filter := publicrequest.FilterRatingSummary{}
	filter.SourceType = input.SourceType
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filter)
		if errMarshal != nil {
			return nil, message.ErrUnmarshalFilterListRatingRequest
		}
	}
	if errValidate := filter.ValidateSourceUID(); errValidate != nil {
		return nil, *errValidate
	}

	ratingSubs, err := s.publicRatingMpRepo.GetPublicRatingSubmissionsGroupBySource(filter)
	if err != nil {
		return nil, message.RecordNotFound
	}
	if len(ratingSubs) <= 0 {
		return results, message.SuccessMsg
	}

	formulaRating, err := s.publicRatingMpRepo.GetRatingFormulaBySourceType(input.SourceType)
	if err != nil {
		return nil, message.RecordNotFound
	}
	if formulaRating == nil {
		return nil, message.RecordNotFound
	}

	// https://it-mkt.atlassian.net/browse/MP-675
	// case product 1: total 10 review, bintang 5 ada 9, bar hijau hampir penuh (9/10 = 90%).
	// case product 2: total 155 review, bintang 4 ada 11, bar hijau nya sedikit (11/155 = 7%)

	// processing calculate detail summary
	for _, ratingSub := range ratingSubs {
		// initiate rating value
		var arrRatingValue []string = global.GetListRatingValueBySourceType(ratingSub.ID.SourceType)

		pRsldr := publicresponse.PublicRatingSummaryListDetailResponse{}
		pRsldr.SourceType = ratingSub.ID.SourceType
		pRsldr.SourceUID = ratingSub.ID.SourceUID
		pRsldr.TotalReviewer = int64(ratingSub.TotalReviewer)
		var totalValue = int64(ratingSub.TotalValue)

		// calculate star
		var arrRatingDetailSummary = populateStarRatingSummary(arrRatingValue, ratingSub.ArrayValue, pRsldr.TotalReviewer)
		pRsldr.RatingSummaryDetail = arrRatingDetailSummary

		// calculate total_value
		sumCountRatingSub := &publicresponse.PublicSumCountRatingSummaryMp{
			Sum:   totalValue,
			Count: pRsldr.TotalReviewer,
		}
		pRsldr.MaximumValue = global.GetMaximumValueBySourceType(input.SourceType)

		ratingSummary, err := calculateRatingMpValue(ratingSub.ID.SourceUID, formulaRating.Formula, sumCountRatingSub)
		if err == nil {
			pRsldr.TotalValue = ratingSummary.TotalValue
			pRsldr.TotalComment = ratingSummary.TotalComment
		}

		results = append(results, pRsldr)
	}

	return results, message.SuccessMsg
}

// swagger:route GET /public/ratings-summary/store-product PublicRating PublicGetRatingSummaryStoreProductRequest
// Get List Rating Summary Store Product
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *publicRatingMpServiceImpl) GetRatingSummaryStoreProduct(ctx context.Context, input publicrequest.PublicGetRatingSummaryStoreProductRequest) ([]publicresponse.RatingSummaryStoreProductNumeric, message.Message) {
	// https://it-mkt.atlassian.net/browse/MP-694
	results := []publicresponse.RatingSummaryStoreProductNumeric{}
	input.MakeDefaultValueIfEmpty()
	var sourceType string = "product"
	filter := publicrequest.FilterRatingSummary{}
	filter.SourceType = sourceType
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filter)
		if errMarshal != nil {
			fmt.Println("errMarshal", errMarshal)
			return nil, message.ErrUnmarshalFilterListRatingRequest
		}
	}
	if errValidate := filter.ValidateStoreUID(); errValidate != nil {
		return nil, *errValidate
	}

	ratingSubs, err := s.publicRatingMpRepo.GetPublicRatingSubmissionsGroupByStoreSource(filter)
	if err != nil {
		return nil, message.RecordNotFound
	}
	if len(ratingSubs) <= 0 {
		return nil, message.SuccessMsg
	}

	formulaRating, err := s.publicRatingMpRepo.GetRatingFormulaBySourceType(sourceType)
	if err != nil {
		return nil, message.RecordNotFound
	}
	if formulaRating == nil {
		return nil, message.RecordNotFound
	}

	// processing calculate  summary
	for _, ratingSub := range ratingSubs {
		result := publicresponse.RatingSummaryStoreProductNumeric{}
		result.StoreUID = ratingSub.ID.StoreUID
		result.TotalReviewer = int64(ratingSub.TotalReviewer)
		result.MaximumValue = global.GetMaximumValueBySourceType(sourceType)
		var totalValue = int64(ratingSub.TotalValue)
		var arrRatingValue []string = global.GetListRatingValueBySourceType(ratingSub.ID.SourceType)

		sumCountRatingSub := &publicresponse.PublicSumCountRatingSummaryMp{
			Sum:   totalValue,
			Count: result.TotalReviewer,
		}
		ratingSummary, err := calculateRatingMpValue(ratingSub.ID.StoreUID, formulaRating.Formula, sumCountRatingSub)
		if err == nil {
			result.TotalValue = ratingSummary.TotalValue
			result.TotalComment = result.TotalReviewer
		}

		// calculate star
		var arrRatingDetailSummary = populateStarRatingSummary(arrRatingValue, ratingSub.ArrayValue, result.TotalReviewer)
		result.RatingSummaryDetail = arrRatingDetailSummary

		results = append(results, result)
	}

	return results, message.SuccessMsg
}

func populateStarRatingSummary(arrRatingValue []string, arrayValue []map[string]int, totalReviewer int64) []publicresponse.PublicRatingSummaryDetailMpResponse {
	var arrRatingDetailSummary []publicresponse.PublicRatingSummaryDetailMpResponse
	for _, arv := range arrRatingValue {
		ratingDetailSummary := publicresponse.PublicRatingSummaryDetailMpResponse{}
		for _, av := range arrayValue {
			ratingDetailSummary.Value = arv

			// assign count value
			if key, isKey := av["key"]; isKey && fmt.Sprint(key) == arv {
				countValue, isCountValue := av["value"]
				if isCountValue {
					ratingDetailSummary.Count = int64(countValue)
				}
			}
		}

		// calculate percentage
		if ratingDetailSummary.Count > 0 {
			percent, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", (float32(ratingDetailSummary.Count)/float32(totalReviewer))*100), 32)
			ratingDetailSummary.Percent = float32(percent)
		}
		arrRatingDetailSummary = append(arrRatingDetailSummary, ratingDetailSummary)
	}

	return arrRatingDetailSummary
}
