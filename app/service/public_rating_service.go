package service

import (
	"encoding/json"
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"
	"strconv"

	"github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PublicRatingService interface {
	GetRatingBySourceTypeAndActor(input request.GetRatingBySourceTypeAndActorRequest) (*response.RatingBySourceTypeAndActorResponse, message.Message)
	CreateRatingSubHelpful(input request.CreateRatingSubHelpfulRequest) message.Message
	GetListRatingSummaryBySourceType(input request.GetPublicListRatingSummaryRequest) ([]response.PublicRatingSummaryResponse, *base.Pagination, message.Message)
	GetListRatingSubmissionBySourceTypeAndUID(input request.GetPublicListRatingSubmissionRequest) ([]response.PublicRatingSubmissionResponse, *base.Pagination, message.Message)
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

// swagger:route GET /api/v1/public/ratings/{source_type}/{source_uid} PublicRating GetRatingBySourceTypeAndActor
// Get Rating By Source Type and Source UID
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *publicRatingServiceImpl) GetRatingBySourceTypeAndActor(input request.GetRatingBySourceTypeAndActorRequest) (*response.RatingBySourceTypeAndActorResponse, message.Message) {
	result := response.RatingBySourceTypeAndActorResponse{}
	// Get Ratings By Type and Actor UID
	ratings, err := s.publicRatingRepo.GetRatingsBySourceTypeAndActor(input.SourceType, input.SourceUID)
	if err != nil {
		return nil, message.FailedMsg
	}
	if len(ratings) == 0 {
		return nil, message.ErrNoData
	}

	result.SourceType = input.SourceType
	result.SourceUID = input.SourceUID
	// Get Rating Type from Rating
	for _, v := range ratings {
		// check rating type exist
		ratingTypeId, err := primitive.ObjectIDFromHex(v.RatingTypeId)
		if err != nil {
			return nil, message.ErrRatingTypeNotExist
		}

		likert, err := s.publicRatingRepo.GetRatingTypeLikertById(ratingTypeId)
		if err != nil {
			return nil, message.FailedMsg
		}
		if likert != nil {
			likertResp := response.MapRatingLikertToRatingNumericResp(*likert)
			result.Ratings = append(result.Ratings, likertResp)
		} else {
			numeric, err := s.publicRatingRepo.GetRatingTypeNumById(ratingTypeId)
			if err != nil {
				return nil, message.FailedMsg
			}
			numericResp := response.MapRatingNumericToRatingNumericResp(*numeric)
			result.Ratings = append(result.Ratings, numericResp)
		}
	}
	return &result, message.SuccessMsg
}

// swagger:route POST /api/v1/public/helpful_rating_submission/ PublicRating ReqRatingSubHelpfulBody
// Create Helpful Rating Submission
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *publicRatingServiceImpl) CreateRatingSubHelpful(input request.CreateRatingSubHelpfulRequest) message.Message {
	// check rating type exist
	ratingSubmissionId, err := primitive.ObjectIDFromHex(input.RatingSubmissionID)
	if err != nil {
		return message.RatingSubmissionNotFound
	}

	ratingSubmission, err1 := s.ratingRepo.GetRatingSubmissionById(ratingSubmissionId)
	if err1 != nil {
		if !errors.Is(err1, mongo.ErrNoDocuments) {
			return message.FailedMsg
		}
	}
	if ratingSubmission == nil {
		return message.ErrRatingTypeNotExist
	}

	_, err2 := s.publicRatingRepo.CreateRatingSubHelpful(input)
	if err2 != nil {
		if mongo.IsDuplicateKeyError(err) {
			return message.ErrDuplicateRatingName
		}
		return message.FailedMsg
	}

	// update like_counter rating submission
	err3 := s.publicRatingRepo.UpdateCounterRatingSubmission(ratingSubmissionId, ratingSubmission.LikeCounter)
	if err3 != nil {
		return message.ErrSaveData
	}
	return message.SuccessMsg
}

// swagger:route GET /api/v1/public/ratings-summary/{source_type} PublicRating GetPublicListRatingSummaryRequest
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
		return nil, nil, message.FailedMsg
	}
	if len(ratings) <= 0 {
		return results, pagination, message.SuccessMsg
	}

	for _, args := range ratings {
		ratingSubs, err := s.publicRatingRepo.GetRatingSubsByRatingId(args.ID.Hex())
		if err != nil {
			return nil, nil, message.FailedMsg
		}
		if len(ratingSubs) <= 0 {
			ratingSummary := response.RatingSubmissionSummary{
				SourceUID:  args.SourceUid,
				TotalValue: 0,
			}
			data := response.PublicRatingSummaryResponse{
				ID:            args.ID,
				Name:          args.Name,
				Description:   args.Description,
				SourceUid:     args.SourceUid,
				SourceType:    args.SourceType,
				RatingType:    args.RatingType,
				RatingTypeId:  args.RatingTypeId,
				RatingSummary: ratingSummary,
			}
			results = append(results, data)
		} else {
			ratingSummary, err := calculateRatingValue(args.SourceUid, ratingSubs)
			if err != nil {
				return []response.PublicRatingSummaryResponse{}, nil, message.ErrFailedToCalculate
			}
			data := response.PublicRatingSummaryResponse{
				ID:            args.ID,
				Name:          args.Name,
				Description:   args.Description,
				SourceUid:     args.SourceUid,
				SourceType:    args.SourceType,
				RatingType:    args.RatingType,
				RatingTypeId:  args.RatingTypeId,
				RatingSummary: ratingSummary,
			}
			results = append(results, data)
		}
	}
	return results, pagination, message.SuccessMsg
}

func calculateRatingValue(sourceUID string, ratingSubs []entity.RatingSubmisson) (response.RatingSubmissionSummary, error) {
	result := response.RatingSubmissionSummary{}
	totalValue := 0
	for _, args := range ratingSubs {
		intVal, err := strconv.Atoi(args.Value)
		if err != nil {
			return result, err
		}
		totalValue = totalValue + intVal
	}
	finalValue := totalValue / len(ratingSubs)
	result.SourceUID = sourceUID
	result.TotalValue = finalValue

	return result, nil
}

// swagger:route GET /api/v1/public/rating-submission/{source_type}/{source_uid} PublicRating GetPublicListRatingSubmissionRequest
// Get Rating Submission By Source Type and Source UID
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *publicRatingServiceImpl) GetListRatingSubmissionBySourceTypeAndUID(input request.GetPublicListRatingSubmissionRequest) ([]response.PublicRatingSubmissionResponse, *base.Pagination, message.Message) {
	results := []response.PublicRatingSubmissionResponse{}
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
		input.Sort = "created at"
	}

	// Get Rating By SourceType and SourceUID
	filterRating := request.FilterRatingSummary{
		SourceType: input.SourceType,
		SourceUid:  input.SourceUID,
	}
	ratings, _, err := s.publicRatingRepo.GetPublicRatingsByParams(input.Limit, input.Page, dir, input.Sort, filterRating)
	if err != nil {
		return nil, nil, message.FailedMsg
	}

	// Get Rating Submission
	filterRatingSubs := request.FilterRatingSubmission{}
	for _, v := range ratings {
		filterRatingSubs.RatingID = append(filterRatingSubs.RatingID, v.ID.Hex())
	}
	ratingSubs, pagination, err := s.publicRatingRepo.GetPublicRatingSubmissions(input.Limit, input.Page, dir, input.Sort, filterRatingSubs)
	if err != nil {
		return nil, nil, message.FailedMsg
	}
	if len(ratings) <= 0 {
		return results, pagination, message.SuccessMsg
	}

	for _, v := range ratingSubs {
		results = append(results, response.PublicRatingSubmissionResponse{
			ID:            v.ID,
			UserID:        v.UserID,
			UserIDLegacy:  v.UserIDLegacy,
			Comment:       v.Comment,
			SourceTransID: v.SourceTransID,
			LikeCounter:   v.LikeCounter,
		})
	}
	return results, pagination, message.SuccessMsg
}
