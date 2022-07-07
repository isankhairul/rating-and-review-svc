package service

import (
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"

	"github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PublicRatingService interface {
	GetRatingBySourceTypeAndActor(input request.GetRatingBySourceTypeAndActorRequest) (*response.RatingBySourceTypeAndActorResponse, message.Message)
}

type publicRatingServiceImpl struct {
	logger           log.Logger
	publicRatingRepo repository.PublicRatingRepository
}

func NewPublicRatingService(
	lg log.Logger,
	prr repository.PublicRatingRepository,
) PublicRatingService {
	return &publicRatingServiceImpl{lg, prr}
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
