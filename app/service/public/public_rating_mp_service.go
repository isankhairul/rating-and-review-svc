package publicservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/log"
	"github.com/spf13/viper"
	"github.com/vjeantet/govaluate"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	publicrequest "go-klikdokter/app/model/request/public"
	"go-klikdokter/app/model/response/public"
	"go-klikdokter/app/repository"
	publicrepository "go-klikdokter/app/repository/public"
	"go-klikdokter/helper/config"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"time"
)

type PublicRatingMpService interface {
	GetListRatingSubmissionBySourceTypeAndUID(input publicrequest.GetPublicListRatingSubmissionRequest) ([]publicresponse.PublicRatingSubmissionMpResponse, *base.Pagination, message.Message)
	GetListRatingSummaryBySourceType(input publicrequest.GetPublicListRatingSummaryRequest) ([]publicresponse.PublicRatingSummaryMpResponse, *base.Pagination, message.Message)
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
	//Set default value
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 50
	}
	if input.Sort == "" {
		input.Sort = "created_at"
	}

	ratings, err := s.publicRatingMpRepo.GetListRatingBySourceTypeAndUID(input.SourceType, input.SourceUID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, message.Message{
				Code:    message.ValidationFailCode,
				Message: "Cannot find rating with params SourceType :" + input.SourceType + ", SourceUid:" + input.SourceUID,
			}
		}
		return nil, nil, message.FailedMsg
	}
	if len(ratings) <= 0 {
		return nil, nil, message.Message{
			Code:    message.ValidationFailCode,
			Message: "Cannot find rating with params SourceType :" + input.SourceType + ", SourceUid:" + input.SourceUID,
		}
	}

	// Unmarshal filter params
	filterRatingSubs := publicrequest.FilterRatingSubmissionMp{}
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filterRatingSubs)
		if errMarshal != nil {
			return nil, nil, message.ErrUnmarshalFilterListRatingRequest
		}
	}
	// Set rating id to filter
	for _, v := range ratings {
		filterRatingSubs.RatingID = append(filterRatingSubs.RatingID, v.ID.Hex())
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
		// Get Rating value
		ratingId, err := primitive.ObjectIDFromHex(v.RatingID)
		if err != nil {
			return nil, nil, message.FailedMsg
		}
		rating, err := s.ratingMpRepo.GetRatingById(ratingId)
		if err != nil {
			return nil, nil, message.ErrRatingNotFound
		}

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

		results = append(results, publicresponse.PublicRatingSubmissionMpResponse{
			ID:            v.ID,
			UserID:        v.UserID,
			UserIDLegacy:  v.UserIDLegacy,
			DisplayName:   displayName,
			Avatar:        v.Avatar,
			Comment:       v.Comment,
			SourceTransID: v.SourceTransID,
			LikeCounter:   v.LikeCounter,
			RatingType:    rating.RatingType,
			Value:         v.Value,
			LikeByMe:      false,
			MediaPath:     v.MediaPath,
			IsWithMedia:   v.IsWithMedia,
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

	ratings, pagination, err := s.publicRatingMpRepo.GetPublicRatingsByParams(input.Limit, input.Page, dir, input.Sort, filter)

	if err != nil {
		return nil, nil, message.RecordNotFound
	}
	if len(ratings) <= 0 {
		return results, pagination, message.SuccessMsg
	}

	for _, args := range ratings {
		ratingTypeId, err := primitive.ObjectIDFromHex(args.RatingTypeId)
		if err != nil {
			return nil, nil, message.FailedMsg
		}
		ratingTypeLikert, err := s.ratingMpRepo.GetRatingTypeLikertByIdAndStatus(ratingTypeId)
		if err != nil {
			if err != mongo.ErrNoDocuments {
				return nil, nil, message.FailedMsg
			}
		}

		if ratingTypeLikert == nil {
			data, err := s.summaryRatingNumeric(args, input.SourceType)
			if err != nil {
				return nil, nil, message.ErrFailedSummaryRatingNumeric
			}
			results = append(results, *data)
		} else {
			data, err := s.summaryRatingLikert(args, *ratingTypeLikert)
			if err != nil {
				return nil, nil, message.ErrFailedSummaryRatingNumeric
			}
			results = append(results, *data)
		}
	}
	return results, pagination, message.SuccessMsg
}

func (s *publicRatingMpServiceImpl) summaryRatingLikert(rating entity.RatingsMpCol, ratingLikert entity.RatingTypesLikertCol) (*publicresponse.PublicRatingSummaryMpResponse, error) {
	likertSummary := publicresponse.RatingSummaryLikert{}
	var myMap map[string]string
	datas, _ := json.Marshal(ratingLikert)
	json.Unmarshal(datas, &myMap)

	likertSummary.SourceUID = rating.SourceUid
	for i := 1; i <= ratingLikert.NumStatements; i++ {
		var statement string
		if i < 10 {
			statement = fmt.Sprint("statement_0", i)
		} else {
			statement = fmt.Sprint("statement_", i)
		}
		data := myMap[statement]
		if data != "" {
			totalCount, err := s.publicRatingMpRepo.CountRatingSubsByRatingIdAndValue(rating.ID.Hex(), strconv.Itoa(i))
			if err != nil {
				return nil, err
			}
			likertObjCount := make(map[string]interface{})
			likertObjCount["seq_id"] = i
			likertObjCount["value"] = data
			likertObjCount["total_reviewer"] = totalCount
			likertSummary.ValueList = append(likertSummary.ValueList, likertObjCount)
		} else {
			return nil, errors.New("invalid statement value")
		}
	}

	result := publicresponse.PublicRatingSummaryMpResponse{
		ID:            rating.ID,
		Name:          rating.Name,
		Description:   rating.Description,
		SourceUid:     rating.SourceUid,
		SourceType:    rating.SourceType,
		RatingType:    rating.RatingType,
		RatingTypeId:  rating.RatingTypeId,
		RatingSummary: likertSummary,
	}
	return &result, nil
}

func (s *publicRatingMpServiceImpl) summaryRatingNumeric(rating entity.RatingsMpCol, sourceType string) (*publicresponse.PublicRatingSummaryMpResponse, error) {
	sumCountRatingSubs, err := s.publicRatingMpRepo.GetSumCountRatingSubsByRatingId(rating.ID.Hex())
	if err != nil {
		return nil, err
	}
	if sumCountRatingSubs == nil {
		return nil, errors.New("data RatingSubmission not found")
	}

	formulaRating, err := s.publicRatingMpRepo.GetRatingFormulaByRatingTypeIdAndSourceType(rating.RatingTypeId, sourceType)
	if err != nil {
		return nil, err
	}
	if formulaRating == nil {
		return nil, errors.New("Formula rating not found, for source_type: " + sourceType)
	}

	if formulaRating.Formula != "" {
		ratingSummary, err := calculateRatingMpValue(rating.SourceUid, formulaRating.Formula, sumCountRatingSubs)
		if err != nil {
			return nil, err
		}

		data := publicresponse.PublicRatingSummaryMpResponse{
			ID:            rating.ID,
			Name:          rating.Name,
			Description:   rating.Description,
			SourceUid:     rating.SourceUid,
			SourceType:    rating.SourceType,
			RatingType:    rating.RatingType,
			RatingTypeId:  rating.RatingTypeId,
			RatingSummary: ratingSummary,
		}
		return &data, nil
	} else {
		return nil, errors.New("formula string is empty")
	}
}

func calculateRatingMpValue(sourceUID, formula string, sumCountRatingSubs *publicresponse.PublicSumCountRatingSummaryMp) (publicresponse.RatingSummaryMpNumeric, error) {
	result := publicresponse.RatingSummaryMpNumeric{}
	result.SourceUID = sourceUID
	result.TotalReviewer = sumCountRatingSubs.Count
	result.TotalValue = 0

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
		totalValue, err := strconv.ParseFloat(fmt.Sprintf("%.1f", finalCalc.(float64)), 64)
		if err != nil {
			return result, err
		}
		result.TotalValue = totalValue
	}

	return result, nil
}
