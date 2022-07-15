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
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util"
	"math"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/vjeantet/govaluate"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PublicRatingService interface {
	GetRatingBySourceTypeAndActor(input request.GetRatingBySourceTypeAndActorRequest) (*response.RatingBySourceTypeAndActorResponse, message.Message)
	CreateRatingSubHelpful(input request.CreateRatingSubHelpfulRequest) message.Message
	GetListRatingSummaryBySourceType(input request.GetPublicListRatingSummaryRequest) ([]response.PublicRatingSummaryResponse, *base.Pagination, message.Message)
	GetListRatingSubmissionBySourceTypeAndUID(input request.GetPublicListRatingSubmissionRequest) ([]response.PublicRatingSubmissionResponse, *base.Pagination, message.Message)
	CreatePublicRatingSubmission(input request.CreateRatingSubmissionRequest) ([]response.PublicCreateRatingSubmissionResponse, message.Message)
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
			if numeric == nil {
				return nil, message.ErrRatingTypeNotExist
			}
			numericResp := response.MapRatingNumericToRatingNumericResp(*numeric)
			result.Ratings = append(result.Ratings, numericResp)
		}
	}
	return &result, message.SuccessMsg
}

// swagger:route POST /api/v1/public/helpful-rating-submission/ PublicRating ReqRatingSubHelpfulBody
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

	ratingSubmission, err := s.ratingRepo.GetRatingSubmissionById(ratingSubmissionId)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return message.FailedMsg
		}
		return message.FailedMsg
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

		result.TotalValue = int(math.Floor(finalCalc.(float64) + 0.4))
	} else {
		result.TotalValue = totalRatingPoint
	}
	result.SourceUID = sourceUID
	result.TotalReviewer = totalUserCount
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
		SourceUid:  []string{input.SourceUID},
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

// swagger:route POST /api/v1/public/rating-submissions/ PublicRating ReqPublicRatingSubmissionBody
// Create Rating Submission Public
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *publicRatingServiceImpl) CreatePublicRatingSubmission(input request.CreateRatingSubmissionRequest) ([]response.PublicCreateRatingSubmissionResponse, message.Message) {
	var saveReq = make([]request.SaveRatingSubmission, 0)
	var empty = ""
	var falseVar = false
	var trueVar = true
	result := []response.PublicCreateRatingSubmissionResponse{}

	// One of the following user_id and user_id_legacy must be filled
	if input.UserID == nil || *input.UserID == "" {
		input.UserID = &empty
	}
	if input.UserIDLegacy == nil || *input.UserIDLegacy == "" {
		input.UserIDLegacy = &empty
	}
	if input.UserID == &empty && input.UserIDLegacy == &empty {
		return result, message.UserUIDRequired
	}

	for _, argRatings := range input.Ratings {
		// Find rating_type_id by rating_id
		objectRatingId, err := primitive.ObjectIDFromHex(argRatings.ID)
		if err != nil {
			return result, message.ErrRatingNotFound
		}
		rating, err := s.ratingRepo.FindRatingByRatingID(objectRatingId)
		if err != nil {
			return result, message.ErrRatingNotFound
		}
		if rating == nil || rating.Status == &falseVar {
			return result, message.ErrRatingNotFound
		}

		// Validate numeric type value
		objectRatingTypeId, err := primitive.ObjectIDFromHex(rating.RatingTypeId)
		if err != nil {
			return result, message.ErrRatingNumericTypeNotFound
		}
		var validateMsg message.Message

		ratingNumericType, err := s.ratingRepo.FindRatingNumericTypeByRatingTypeID(objectRatingTypeId)
		if err != nil || ratingNumericType.Status == &falseVar {
			validateMsg = message.ErrRatingNumericTypeNotFound
		} else {
			value, er := strconv.ParseFloat(*argRatings.Value, 64)
			if er != nil {
				return result, message.ErrValueFormatForNumericType
			}
			validateMsg = util.ValidateTypeNumeric(ratingNumericType, value)
		}
		if validateMsg.Code == message.ValidationFailCode {
			return result, validateMsg
		}

		// Validate numeric type value
		if validateMsg == message.ErrRatingNumericTypeNotFound {
			ratingTypeLikert, err := s.ratingRepo.GetRatingTypeLikertByIdAndStatus(objectRatingTypeId)
			if err != nil {
				return result, message.ErrLikertTypeNotFound
			}
			if ratingTypeLikert == nil {
				return result, message.ErrLikertTypeNotFound
			}
			strValue := strings.Split(*argRatings.Value, ",")
			if validateErr, validList := util.ValidateLikertType(ratingTypeLikert, strValue); validateErr != nil {
				return result, message.Message{
					Code:    message.ValidationFailCode,
					Message: "value must be integer and include in " + fmt.Sprintf("%v", validList),
				}
			}
		}

		// A submission with a combination of either (rating_id and user_id) OR (rating_id and user_id_legacy) is allowed once
		if input.UserID != &empty && input.UserIDLegacy != &empty {
			ratingSubmission, er := s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(input.UserIDLegacy, argRatings.ID, input.SourceTransID)
			val := util.ValidateUserIdAndUserIdLegacy(input, rating.ID.Hex(), input.UserID, input.UserIDLegacy, ratingSubmission, er)
			if val {
				return result, message.UserRated
			}
			if !val {
				ratingSubmission, er = s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(input.UserID, argRatings.ID, input.SourceTransID)
				valL := util.ValidateUserIdAndUserIdLegacy(input, rating.ID.Hex(), input.UserID, input.UserIDLegacy, ratingSubmission, er)
				if valL {
					return result, message.UserRated

				}
			}
		}

		if input.UserID == &empty {
			ratingSubmission, er := s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(input.UserIDLegacy, argRatings.ID, input.SourceTransID)
			if val := util.ValidateUserIdAndUserIdLegacy(input, rating.ID.Hex(), input.UserID, input.UserIDLegacy, ratingSubmission, er); val == trueVar {
				return result, message.UserRated
			}
		}
		if input.UserIDLegacy == &empty {
			ratingSubmission, er := s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(input.UserID, argRatings.ID, input.SourceTransID)
			if val := util.ValidateUserIdAndUserIdLegacy(input, rating.ID.Hex(), input.UserID, input.UserIDLegacy, ratingSubmission, er); val == trueVar {
				return result, message.UserRated
			}
		}

		//The maximum length of user_agent allowed is 200 characters. Crop at 197 characters with triple dots (...) at the end.
		if len(strings.TrimSpace(input.UserAgent)) > 200 {
			return result, message.UserAgentTooLong
		}
		if isExisted := isIdExisted(saveReq, argRatings.ID); isExisted == falseVar {
			return result, message.ErrCannotSameRatingId
		}

		if ratingNumericType == nil {
			saveReq = append(saveReq, request.SaveRatingSubmission{
				RatingID:      argRatings.ID,
				Value:         argRatings.Value,
				UserID:        input.UserID,
				UserIDLegacy:  input.UserIDLegacy,
				Comment:       "",
				IPAddress:     input.IPAddress,
				UserAgent:     input.UserAgent,
				SourceTransID: input.SourceTransID,
				UserPlatform:  input.UserPlatform,
			})
		} else {
			saveReq = append(saveReq, request.SaveRatingSubmission{
				RatingID:      argRatings.ID,
				Value:         argRatings.Value,
				UserID:        input.UserID,
				UserIDLegacy:  input.UserIDLegacy,
				Comment:       input.Comment,
				IPAddress:     input.IPAddress,
				UserAgent:     input.UserAgent,
				SourceTransID: input.SourceTransID,
				UserPlatform:  input.UserPlatform,
			})
		}
	}
	if len(saveReq) <= 0 {
		return result, message.ErrTypeNotFound
	}
	ratingSubs, err := s.publicRatingRepo.CreatePublicRatingSubmission(saveReq)
	if err != nil {
		return result, message.ErrSaveData
	}

	for _, arg := range ratingSubs {
		data := response.PublicCreateRatingSubmissionResponse{}
		ratingSub, err := s.ratingRepo.GetRatingSubmissionById(arg.ID)
		if err != nil {
			return result, message.FailedMsg
		}
		ratingID, err := primitive.ObjectIDFromHex(ratingSub.RatingID)
		if err != nil {
			return result, message.FailedMsg
		}
		rating, err := s.ratingRepo.GetRatingById(ratingID)
		if err != nil {
			return result, message.FailedMsg
		}

		data.ID = ratingSub.ID
		data.RatingID = ratingSub.RatingID
		data.RatingDescription = *rating.Description
		data.Value = ratingSub.Value

		result = append(result, data)
	}
	return result, message.SuccessMsg
}
