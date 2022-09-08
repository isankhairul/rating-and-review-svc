package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/config"
	"go-klikdokter/helper/message"
	"math"
	"strconv"
	"time"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
	"github.com/vjeantet/govaluate"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PublicRatingService interface {
	GetListRatingSubmissionBySourceTypeAndUID(input request.GetPublicListRatingSubmissionRequest) ([]response.PublicRatingSubmissionResponse, *base.Pagination, message.Message)
	GetListRatingSummaryBySourceType(input request.GetPublicListRatingSummaryRequest) ([]response.PublicRatingSummaryResponse, *base.Pagination, message.Message)
}

type publicRatingServiceImpl struct {
	logger           log.Logger
	ratingRepo       repository.RatingRepository
	publicRatingRepo repository.PublicRatingRepository
}

func NewPublicRatingService(
	lg log.Logger,
	rr repository.RatingRepository,
	prr repository.PublicRatingRepository,
) PublicRatingService {
	return &publicRatingServiceImpl{lg, rr, prr}
}

// swagger:route GET /public/ratings-summary/{source_type} PublicRating GetPublicListRatingSummaryRequest
// Get Rating Summary By Source Type
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *publicRatingServiceImpl) GetListRatingSummaryBySourceType(input request.GetPublicListRatingSummaryRequest) ([]response.PublicRatingSummaryResponse, *base.Pagination, message.Message) {
	results := []response.PublicRatingSummaryResponse{}
	input.MakeDefaultValueIfEmpty()
	var dir int
	if input.Dir == "asc" {
		dir = 1
	} else {
		dir = -1
	}
	filter := request.FilterRatingSummary{}
	filter.SourceType = input.SourceType
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filter)
		if errMarshal != nil {
			return nil, nil, message.ErrUnmarshalFilterListRatingRequest
		}
	}

	ratings, pagination, err := s.publicRatingRepo.GetPublicRatingsByParams(input.Limit, input.Page, dir, input.Sort, filter)
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
		ratingTypeLikert, err := s.ratingRepo.GetRatingTypeLikertByIdAndStatus(ratingTypeId)
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

func (s *publicRatingServiceImpl) summaryRatingLikert(rating entity.RatingsCol, ratingLikert entity.RatingTypesLikertCol) (*response.PublicRatingSummaryResponse, error) {
	likertSummary := response.RatingSummaryLikert{}
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
			totalCount, err := s.publicRatingRepo.CountRatingSubsByRatingIdAndValue(rating.ID.Hex(), strconv.Itoa(i))
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

	result := response.PublicRatingSummaryResponse{
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

func (s *publicRatingServiceImpl) summaryRatingNumeric(rating entity.RatingsCol, sourceType string) (*response.PublicRatingSummaryResponse, error) {
	ratingSubs, err := s.publicRatingRepo.GetRatingSubsByRatingId(rating.ID.Hex())
	if err != nil {
		return nil, err
	}

	formulaRating, err := s.publicRatingRepo.GetRatingFormulaByRatingTypeIdAndSourceType(rating.RatingTypeId, sourceType)
	if err != nil {
		return nil, err
	}
	if formulaRating.Formula != "" {
		ratingSummary, err := calculateRatingValue(rating.SourceUid, formulaRating.Formula, ratingSubs)
		if err != nil {
			return nil, err
		}

		data := response.PublicRatingSummaryResponse{
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

func calculateRatingValue(sourceUID, formula string, ratingSubs []entity.RatingSubmisson) (response.RatingSummaryNumeric, error) {
	result := response.RatingSummaryNumeric{}
	totalRatingPoint := 0
	totalUserCount := len(ratingSubs)

	if totalUserCount > 0 {
		// Get total rating point
		for _, args := range ratingSubs {
			intVal, err := strconv.Atoi(args.Value)
			if err != nil {
				return result, err
			}
			totalRatingPoint = totalRatingPoint + intVal
		}
	}

	if formula != "" {
		expression, err := govaluate.NewEvaluableExpression(formula)
		if err != nil {
			return result, err
		}
		parameters := make(map[string]interface{}, 8)
		parameters["total_rating_point"] = totalRatingPoint
		parameters["total_user_count"] = totalUserCount

		finalCalc, err := expression.Evaluate(parameters)
		if err != nil {
			return result, err
		}

		result.TotalValue = int(math.Floor(finalCalc.(float64) + 0.5))
	} else {
		result.TotalValue = totalRatingPoint
	}
	result.SourceUID = sourceUID
	result.TotalReviewer = totalUserCount
	return result, nil
}

// swagger:route GET /public/rating-submissions/{source_type}/{source_uid} PublicRating GetPublicListRatingSubmissionRequest
// Get Rating Submission By Source Type and Source UID
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *publicRatingServiceImpl) GetListRatingSubmissionBySourceTypeAndUID(input request.GetPublicListRatingSubmissionRequest) ([]response.PublicRatingSubmissionResponse, *base.Pagination, message.Message) {
	results := []response.PublicRatingSubmissionResponse{}
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

	ratings, err := s.publicRatingRepo.GetListRatingBySourceTypeAndUID(input.SourceType, input.SourceUID)
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
	filterRatingSubs := request.FilterRatingSubmission{}
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
	ratingSubs, pagination, err := s.publicRatingRepo.GetPublicRatingSubmissions(input.Limit, input.Page, dir, input.Sort, filterRatingSubs)
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
		rating, err := s.ratingRepo.GetRatingById(ratingId)
		if err != nil {
			return nil, nil, message.ErrRatingNotFound
		}

		if v.Avatar == "" {
			v.Avatar = avatarDefault
		}
		Loc, _ := time.LoadLocation(timezone)

		results = append(results, response.PublicRatingSubmissionResponse{
			ID:            v.ID,
			UserID:        v.UserID,
			UserIDLegacy:  v.UserIDLegacy,
			DisplayName:   *v.DisplayName,
			Avatar:        v.Avatar,
			Comment:       v.Comment,
			SourceTransID: v.SourceTransID,
			LikeCounter:   v.LikeCounter,
			RatingType:    rating.RatingType,
			Value:         v.Value,
			LikeByMe:      false,
			CreatedAt:     v.CreatedAt.In(Loc),
		})
	}
	return results, pagination, message.SuccessMsg
}
