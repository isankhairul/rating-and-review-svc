package publicservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-klikdokter/app/model/base"
	publicrequest "go-klikdokter/app/model/request/public"
	publicresponse "go-klikdokter/app/model/response/public"
	"go-klikdokter/app/repository"
	publicrepository "go-klikdokter/app/repository/public"
	"go-klikdokter/helper/config"
	"go-klikdokter/helper/message"
	"go-klikdokter/helper/thumbor"
	"go-klikdokter/pkg/util"
	"time"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
	"github.com/vjeantet/govaluate"
	"go.mongodb.org/mongo-driver/mongo"
)

type PublicRatingMpService interface {
	GetListRatingSubmissionBySourceTypeAndUID(input publicrequest.GetPublicListRatingSubmissionRequest) ([]publicresponse.PublicRatingSubmissionMpResponse, *base.Pagination, message.Message)
	GetListRatingSummaryBySourceType(input publicrequest.GetPublicListRatingSummaryRequest) ([]publicresponse.PublicRatingSummaryMpResponse, *base.Pagination, message.Message)
	GetListRatingSubmissionByID(ctx context.Context, input publicrequest.GetPublicListRatingSubmissionByIDRequest) ([]publicresponse.PublicRatingSubmissionMpResponse, *base.Pagination, message.Message, interface{})
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
		if v.MediaPath == nil {
			v.MediaPath = []string{}
		}

		// create thumbor response
		mediaImages := []string{}
		for _, value := range v.MediaPath {
			mediaImages = append(mediaImages, thumbor.GetNewThumborImagesOriginal(value))
		}

		results = append(results, publicresponse.PublicRatingSubmissionMpResponse{
			ID:            v.ID,
			SourceType:    v.SourceType,
			SourceUID:     v.SourceUID,
			UserID:        v.UserID,
			UserIDLegacy:  v.UserIDLegacy,
			DisplayName:   displayName,
			Avatar:        v.Avatar,
			Comment:       v.Comment,
			SourceTransID: v.SourceTransID,
			LikeCounter:   v.LikeCounter,
			Value:         v.Value,
			LikeByMe:      false,
			MediaPath:     v.MediaPath,
			IsWithMedia:   v.IsWithMedia,
			MediaImages:   mediaImages,
			CreatedAt:     v.CreatedAt.In(Loc),
		})

	}
	return results, pagination, message.SuccessMsg
}

func (s *publicRatingMpServiceImpl) GetListRatingSummaryBySourceType(input publicrequest.GetPublicListRatingSummaryRequest) ([]publicresponse.PublicRatingSummaryMpResponse, *base.Pagination, message.Message) {
	results := []publicresponse.PublicRatingSummaryMpResponse{}
	input.MakeDefaultValueIfEmpty()
	var dir int
	if input.Dir == "asc" {
		dir = 1
	} else {
		dir = -1
	}
	filter := publicrequest.FilterRatingSummary{}
	filter.SourceType = input.SourceType
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filter)
		if errMarshal != nil {
			return nil, nil, message.ErrUnmarshalFilterListRatingRequest
		}
	}
	if len(filter.SourceUid) == 0 {
		return nil, nil, message.ErrSourceUidRequire
	}

	ratingSub, pagination, err := s.publicRatingMpRepo.GetPublicRatingSubmissionsGroupBySource(input.Limit, input.Page, dir, input.Sort, filter)
	if err != nil {
		return nil, nil, message.RecordNotFound
	}
	if len(ratingSub) <= 0 {
		return results, pagination, message.SuccessMsg
	}

	for _, args := range ratingSub {
		data, err := s.summaryRatingNumeric(args.ID.SourceUID, input.SourceType)
		if err != nil {
			return nil, nil, message.ErrFailedSummaryRatingNumeric
		}
		results = append(results, *data)
	}
	return results, pagination, message.SuccessMsg
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
		input.Limit = 10
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
		if v.IsAnonymous {
			displayName = util.MaskDisplayName(*v.DisplayName)
		} else {
			displayName = *v.DisplayName
		}

		// create thumbor response
		mediaImages := []string{}
		for _, value := range v.MediaPath {
			mediaImages = append(mediaImages, thumbor.GetNewThumborImagesOriginal(value))
		}
		if v.MediaPath == nil {
			v.MediaPath = []string{}
		}
		result = append(result, publicresponse.PublicRatingSubmissionMpResponse{
			ID:            v.ID,
			SourceType:    v.SourceType,
			SourceUID:     v.SourceUID,
			UserID:        v.UserID,
			UserIDLegacy:  v.UserIDLegacy,
			DisplayName:   displayName,
			Avatar:        v.Avatar,
			Comment:       v.Comment,
			SourceTransID: v.SourceTransID,
			LikeCounter:   v.LikeCounter,
			Value:         v.Value,
			LikeByMe:      false,
			MediaPath:     nil,
			IsWithMedia:   v.IsWithMedia,
			MediaImages:   mediaImages,
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

		// totalValue, err := strconv.ParseFloat(fmt.Sprintf("%.1f", finalCalc.(float64)), 64)
		// if err != nil {
		//	return result, err
		// }
		result.TotalValue = fmt.Sprintf("%.1f", finalCalc)
	}

	return result, nil
}
