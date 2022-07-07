package service

import (
	"errors"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"

	"github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PublicRatingService interface {
	GetRatingBySourceTypeAndActor(input request.GetRatingBySourceTypeAndActorRequest) (*response.RatingBySourceTypeAndActorResponse, message.Message)
	CreateRatingSubHelpful(input request.CreateRatingSubHelpfulRequest) message.Message
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
