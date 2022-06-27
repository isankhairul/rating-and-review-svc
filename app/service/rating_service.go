package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/log"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"strconv"
	"strings"
)

type RatingService interface {
	// Rating type num
	CreateRatingTypeNum(input request.CreateRatingTypeNumRequest) (*entity.RatingTypesNumCol, message.Message)
	UpdateRatingTypeNum(input request.EditRatingTypeNumRequest) message.Message
	GetRatingTypeNumById(input request.GetRatingTypeNumRequest) (*entity.RatingTypesNumCol, message.Message)
	DeleteRatingTypeNumById(input request.GetRatingTypeNumRequest) message.Message
	GetRatingTypeNums(input request.GetRatingTypeNumsRequest) ([]entity.RatingTypesNumCol, *base.Pagination, message.Message)

	// Rating submission
	CreateRatingSubmission(input request.CreateRatingSubmissionRequest) message.Message
	UpdateRatingSubmission(input request.UpdateRatingSubmissonRequest) message.Message
	GetRatingSubmission(id string) (*response.RatingSubmissonResponse, message.Message)
	GetListRatingSubmissions(input request.ListRatingSubmissionRequest) ([]response.RatingSubmissonResponse, *base.Pagination, message.Message)
	DeleteRatingSubmission(id string) message.Message

	// Rating type likert
	CreateRatingTypeLikert(input request.SaveRatingTypeLikertRequest) message.Message
	GetRatingTypeLikertById(input request.GetRatingTypeLikertRequest) (*entity.RatingTypesLikertCol, message.Message)
	UpdateRatingTypeLikert(input request.SaveRatingTypeLikertRequest) message.Message
	DeleteRatingTypeLikertById(input request.GetRatingTypeLikertRequest) message.Message
	GetRatingTypeLikerts(input request.GetRatingTypeLikertsRequest) ([]entity.RatingTypesLikertCol, *base.Pagination, message.Message)

	// Rating
	CreateRating(input request.SaveRatingRequest) (*entity.RatingsCol, message.Message)
	GetRatingById(id string) (*entity.RatingsCol, message.Message)
	UpdateRating(input request.UpdateRatingRequest) message.Message
	DeleteRating(id string) message.Message
	GetListRatings(input request.GetListRatingsRequest) ([]entity.RatingsCol, *base.Pagination, message.Message)
	GetListRatingSummary(input request.GetListRatingSummaryRequest) ([]response.RatingSummaryResponse, message.Message)
}

type ratingServiceImpl struct {
	logger          log.Logger
	ratingRepo      repository.RatingRepository
	medicalFacility util.MedicalFacilitySvc
}

func NewRatingService(
	lg log.Logger,
	rr repository.RatingRepository,
	mf util.MedicalFacilitySvc,
) RatingService {
	return &ratingServiceImpl{lg, rr, mf}
}

// swagger:route POST /api/v1/rating-types-numeric/ RatingTypeNum createRatingTypeNum
// Create Numeric Rating Types
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) CreateRatingTypeNum(input request.CreateRatingTypeNumRequest) (*entity.RatingTypesNumCol, message.Message) {
	if *input.Scale < 0 || *input.Scale > 2 {
		return nil, message.ErrScaleValueReq
	}
	if *input.MaxScore <= *input.MinScore {
		return nil, message.ErrMaxMin
	}
	check := true
	if input.Status == nil {
		input.Status = &check
	}

	interval := util.ValidInterval(*input.MinScore, *input.MaxScore, *input.Scale)
	if input.Intervals != interval {
		return nil, message.Message{
			Code:    message.ValidationFailCode,
			Message: "interval must be " + strconv.Itoa(interval),
		}
	}

	result, err := s.ratingRepo.CreateRatingTypeNum(input)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, message.ErrDuplicateType
		}
		return nil, message.FailedMsg
	}
	return result, message.SuccessMsg
}

// swagger:route GET /api/v1/rating-types-numeric/{id} RatingTypeNum getRatingTypeNumById
// Get Numeric Rating Types by ID
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) GetRatingTypeNumById(input request.GetRatingTypeNumRequest) (*entity.RatingTypesNumCol, message.Message) {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return nil, message.ErrNoData
	}
	result, err := s.ratingRepo.GetRatingTypeNumById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.ErrNoData
		}
		return nil, message.FailedMsg
	}
	return result, message.SuccessMsg
}

func checkConditionUpdateRatingTypeNum(s *ratingServiceImpl, input request.EditRatingTypeNumRequest) message.Message {
	rating, err := s.ratingRepo.GetRatingByType(input.Id)
	if err != nil && err != mongo.ErrNoDocuments {
		return message.FailedMsg
	}
	if rating != nil {
		msg := util.ValidInputUpdateRatingTypeNumRated(input)
		if msg != message.SuccessMsg {
			return msg
		}
		submissison, err := s.ratingRepo.GetRatingSubmissionByRatingId(rating.ID.Hex())
		if err != nil && err != mongo.ErrNoDocuments {
			return message.FailedMsg
		}
		if submissison != nil {
			msg := util.ValidInputUpdateRatingTypeNumSubmission(input)
			if msg != message.SuccessMsg {
				return msg
			}
		}
	} else {
		interval := util.ValidInterval(*input.MinScore, *input.MaxScore, *input.Scale)
		if *input.Intervals != interval {
			return message.Message{
				Code:    message.ValidationFailCode,
				Message: "interval must be " + strconv.Itoa(interval),
			}
		}
	}
	return message.SuccessMsg
}

// swagger:route PUT /api/v1/rating-types-numeric/{id} RatingTypeNum updateRatingTypeNum
// Update Numeric Rating Type
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) UpdateRatingTypeNum(input request.EditRatingTypeNumRequest) message.Message {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return message.ErrNoData
	}

	msg := checkConditionUpdateRatingTypeNum(s, input)
	if msg != message.SuccessMsg {
		return msg
	}

	err = s.ratingRepo.UpdateRatingTypeNum(objectId, input)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return message.ErrNoData
		}
		if mongo.IsDuplicateKeyError(err) {
			return message.ErrDuplicateType
		}
		return message.FailedMsg
	}
	return message.SuccessMsg
}

// swagger:route Delete /api/v1/rating-types-numeric/{id} RatingTypeNum deleteRatingTypeNum
// Delete Numeric Rating Type
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) DeleteRatingTypeNumById(input request.GetRatingTypeNumRequest) message.Message {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return message.ErrNoData
	}
	rating, err := s.ratingRepo.GetRatingByType(input.Id)
	if err != nil && err != mongo.ErrNoDocuments {
		return message.FailedMsg
	}
	if rating != nil {
		submissison, err := s.ratingRepo.GetRatingSubmissionByRatingId(rating.ID.Hex())
		if err != nil && err != mongo.ErrNoDocuments {
			return message.FailedMsg
		}
		if submissison != nil {
			return message.ErrThisRatingTypeIsInUse
		}
	}
	err = s.ratingRepo.DeleteRatingTypeNum(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return message.ErrNoData
		}
		return message.FailedMsg
	}
	return message.SuccessMsg
}

// swagger:route GET /api/v1/rating-types-numeric RatingTypeNum getRatingTypeNums
// Get Numeric Rating Types Listing
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) GetRatingTypeNums(input request.GetRatingTypeNumsRequest) ([]entity.RatingTypesNumCol, *base.Pagination, message.Message) {
	var dir interface{}
	sort := "updated_at"
	if input.Sort != "" {
		sort = input.Sort
	}
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

	filter := request.Filter{}
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filter)
		if errMarshal != nil {
			return nil, nil, message.ErrUnmarshalRequest
		}
	}
	ratingTypeNums, pagination, err := s.ratingRepo.GetRatingTypeNums(filter, input.Page, input.Limit, sort, dir)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, message.ErrNoData
		}
		return nil, nil, message.FailedMsg
	}
	results := make([]entity.RatingTypesNumCol, 0)
	if len(ratingTypeNums) == 0 {
		return results, pagination, message.SuccessMsg
	}
	results = ratingTypeNums

	return results, pagination, message.SuccessMsg
}

// swagger:route POST /api/v1/rating-submissions/ RatingSubmission ReqRatingSubmissonBody
// Create Rating Submission
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) CreateRatingSubmission(input request.CreateRatingSubmissionRequest) message.Message {

	// One of the following user_id and user_id_legacy must be filled.
	if input.UserID == nil && input.UserIDLegacy == nil {
		return message.UserUIDRequired
	}

	var saveReq = make([]request.SaveRatingSubmission, 0)

	for _, argRatings := range input.Ratings {

		//The value must be valid according to requirements of rating type
		// Find rating_type_id by rating_id
		objectRatingId, err := primitive.ObjectIDFromHex(argRatings.ID)
		if err != nil {
			return message.ErrRatingNotFound
		}
		rating, err := s.ratingRepo.FindRatingByRatingID(objectRatingId)
		if err != nil {
			return message.ErrRatingNotFound
		}
		if rating == nil || *rating.Status == false {
			return message.ErrRatingNotFound
		}

		// Validate value
		objectRatingTypeId, err := primitive.ObjectIDFromHex(rating.RatingTypeId)
		if err != nil {
			return message.ErrRatingNumericTypeNotFound
		}
		var validateMsg message.Message

		// Numeric type
		ratingNumericType, err := s.ratingRepo.FindRatingNumericTypeByRatingTypeID(objectRatingTypeId)
		if ratingNumericType == nil || *ratingNumericType.Status == false || err != nil {
			validateMsg = message.ErrRatingNumericTypeNotFound
		} else {
			value, er := strconv.ParseFloat(*argRatings.Value, 64)
			if er != nil {
				return message.ErrValueFormat
			}
			validateMsg = util.ValidateTypeNumeric(ratingNumericType, value)
		}

		if validateMsg.Code == message.ValidationFailCode {
			return validateMsg
		}

		// Likert type
		if validateMsg == message.ErrRatingNumericTypeNotFound {
			ratingTypeLikert, err := s.ratingRepo.GetRatingTypeLikertByIdAndStatus(objectRatingTypeId)
			if err != nil {
				if errors.Is(mongo.ErrNoDocuments, err) {
					return message.ErrTypeNotFound
				}
				return message.FailedMsg
			}
			if ratingTypeLikert == nil {
				return message.ErrTypeNotFound
			}
			strValue := strings.Split(*argRatings.Value, ",")
			if validateErr, validList := util.ValidateLikertType(ratingTypeLikert, strValue); validateErr != nil {
				return message.Message{
					Code:    message.ValidationFailCode,
					Message: "value must be integer and include in " + fmt.Sprintf("%v", validList),
				}
			}
		}

		// A submission with a combination of either (rating_id and user_id) OR (rating_id and user_id_legacy) is allowed once
		if input.UserID != nil && input.UserIDLegacy != nil && *input.UserID != "" && *input.UserIDLegacy != "" {
			ratingSubmission, err := s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(input.UserIDLegacy, argRatings.ID)
			if ratingSubmission != nil && err == nil {
				return message.UserRated
			}
			ratingSubmission, err = s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(input.UserID, argRatings.ID)
			if ratingSubmission != nil && err == nil {
				return message.UserRated
			}
		}
		if input.UserID == nil || *input.UserID == "" {
			ratingSubmission, err := s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(input.UserIDLegacy, argRatings.ID)
			if ratingSubmission != nil && err == nil {
				return message.UserRated
			}
		}
		if input.UserIDLegacy == nil || *input.UserIDLegacy == "" {
			ratingSubmission, err := s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(input.UserID, argRatings.ID)
			if ratingSubmission != nil && err == nil {
				return message.UserRated
			}
		}

		// A submission with a combination of either (rating_id and user_id) OR (rating_id and user_id_legacy) is allowed once
		if input.UserID != nil && input.UserIDLegacy != nil && *input.UserID != "" && *input.UserIDLegacy != "" {
			ratingSubmission, err := s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(input.UserIDLegacy, argRatings.ID)
			if ratingSubmission != nil && err == nil {
				return message.UserRated
			}
			ratingSubmission, err = s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(input.UserID, argRatings.ID)
			if ratingSubmission != nil && err == nil {
				return message.UserRated
			}
		}
		if input.UserID == nil || *input.UserID == "" {
			ratingSubmission, err := s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(input.UserIDLegacy, argRatings.ID)
			if ratingSubmission != nil && err == nil {
				return message.UserRated
			}
		}
		if input.UserIDLegacy == nil || *input.UserIDLegacy == "" {
			ratingSubmission, err := s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(input.UserID, argRatings.ID)
			if ratingSubmission != nil && err == nil {
				return message.UserRated
			}
		}

		//The maximum length of user_agent allowed is 200 characters. Crop at 197 characters with triple dots (...) at the end.
		if len(strings.TrimSpace(input.UserAgent)) > 200 {
			return message.UserAgentTooLong
		}

		if isExisted := isIdExisted(saveReq, argRatings.ID); isExisted == false {
			return message.ErrCannotSameRatingId
		}

		if ratingNumericType == nil {
			saveReq = append(saveReq, request.SaveRatingSubmission{
				RatingID:     argRatings.ID,
				Value:        argRatings.Value,
				UserID:       input.UserID,
				UserIDLegacy: input.UserIDLegacy,
				Comment:      "",
				IPAddress:    input.IPAddress,
				UserAgent:    input.UserAgent,
			})
		} else {
			saveReq = append(saveReq, request.SaveRatingSubmission{
				RatingID:     argRatings.ID,
				Value:        argRatings.Value,
				UserID:       input.UserID,
				UserIDLegacy: input.UserIDLegacy,
				Comment:      input.Comment,
				IPAddress:    input.IPAddress,
				UserAgent:    input.UserAgent,
			})
		}
	}
	if len(saveReq) < 0 {
		return message.ErrTypeNotFound
	}
	_, err := s.ratingRepo.CreateRatingSubmission(saveReq)
	if err != nil {
		return message.ErrSaveData
	}
	return message.SuccessMsg
}

func isIdExisted(saveReq []request.SaveRatingSubmission, ratingId string) bool {
	for _, args := range saveReq {
		if ratingId == args.RatingID {
			return false
		}
	}
	return true
}

// swagger:route PUT /api/v1/rating-submissions/{id} RatingSubmission ReqUpdateRatingSubmissonBody
// Update Rating Submission
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) UpdateRatingSubmission(input request.UpdateRatingSubmissonRequest) message.Message {

	// One of the following user_id and user_id_legacy must be filled.
	if input.UserID == nil && input.UserIDLegacy == nil {
		return message.UserUIDRequired
	}

	// find Rating submission
	objectRatingSubmissionId, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return message.RatingSubmissionNotFound
	}
	ratingSubmission, err := s.ratingRepo.GetRatingSubmissionById(objectRatingSubmissionId)
	if err != nil || ratingSubmission == nil {
		return message.RatingSubmissionNotFound
	}

	//The value must be valid according to requirements of rating type
	// Find rating_type_id by rating_id
	objectRatingId, err := primitive.ObjectIDFromHex(input.RatingID)
	if err != nil {
		return message.ErrRatingNotFound
	}
	rating, err := s.ratingRepo.FindRatingByRatingID(objectRatingId)
	if err != nil {
		return message.ErrRatingNotFound
	}
	if rating == nil || *rating.Status == false {
		return message.ErrRatingNotFound
	}

	// Validate value
	objectRatingTypeId, err := primitive.ObjectIDFromHex(rating.RatingTypeId)
	if err != nil {
		return message.ErrTypeNotFound
	}
	var validateMsg message.Message

	// Numeric type
	ratingNumericType, err := s.ratingRepo.FindRatingNumericTypeByRatingTypeID(objectRatingTypeId)
	if ratingNumericType == nil || *ratingNumericType.Status == false || err != nil {
		validateMsg = message.ErrRatingNumericTypeNotFound
	} else {
		value, er := strconv.ParseFloat(*input.Value, 64)
		if er != nil {
			return message.ErrValueFormat
		}
		validateMsg = util.ValidateTypeNumeric(ratingNumericType, value)
	}

	if validateMsg.Code == message.ValidationFailCode {
		return validateMsg
	}

	// Likert type
	if validateMsg == message.ErrRatingNumericTypeNotFound {
		ratingTypeLikert, err := s.ratingRepo.GetRatingTypeLikertByIdAndStatus(objectRatingTypeId)
		if err != nil {
			if errors.Is(mongo.ErrNoDocuments, err) {
				return message.ErrTypeNotFound
			}
			return message.FailedMsg
		}
		if ratingTypeLikert == nil {
			return message.ErrTypeNotFound
		}
		strValue := strings.Split(*input.Value, ",")
		if validateErr, validList := util.ValidateLikertType(ratingTypeLikert, strValue); validateErr != nil {
			return message.Message{
				Code:    message.ValidationFailCode,
				Message: "value must be integer and include in " + fmt.Sprintf("%v", validList),
			}
		}
		input.Comment = ""
	}

	// A submission with a combination of either (rating_id and user_id) OR (rating_id and user_id_legacy) is allowed once
	if input.UserID != nil && input.UserIDLegacy != nil && *input.UserID != "" && *input.UserIDLegacy != "" {
		ratingSubmission, err = s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(input.UserIDLegacy, input.ID)
		if ratingSubmission != nil && err == nil {
			return message.UserRated
		}
		ratingSubmission, err = s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(input.UserID, input.ID)
		if ratingSubmission != nil && err == nil {
			return message.UserRated
		}
	}
	if input.UserID == nil || *input.UserID == "" {
		ratingSubmission, err = s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(input.UserIDLegacy, input.ID)
		if ratingSubmission != nil && err == nil {
			return message.UserRated
		}
	}
	if input.UserIDLegacy == nil || *input.UserIDLegacy == "" {
		ratingSubmission, err = s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(input.UserID, input.ID)
		if ratingSubmission != nil && err == nil {
			return message.UserRated
		}
	}

	objecUpdateId, _ := primitive.ObjectIDFromHex(input.ID)
	errC := s.ratingRepo.UpdateRatingSubmission(input, objecUpdateId)
	if errC != nil {
		return message.ErrSaveData
	}
	return message.SuccessMsg

}

// swagger:route DELETE /api/v1/rating-submissions/{id} RatingSubmission ReqDeleteRatingSubmissionById
// Delete Rating Submission
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) DeleteRatingSubmission(id string) message.Message {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return message.ErrDataNotFound
	}
	err = s.ratingRepo.DeleteSubmission(objectId)
	if err != nil {
		if err.Error() == "user not found" {
			return message.ErrNoData
		}
		return message.FailedMsg
	}
	return message.SuccessMsg
}

// swagger:route GET /api/v1/rating-submissions/{id} RatingSubmission ReqRatingSubmissionById
// Get Rating Submission By Id
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) GetRatingSubmission(id string) (*response.RatingSubmissonResponse, message.Message) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, message.ErrDataNotFound
	}
	get, err := s.ratingRepo.GetRatingSubmissionById(objectId)
	if err != nil {
		fmt.Println("err.Error()", err.Error())
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.ErrNoData
		}
		return nil, message.FailedMsg
	}
	var result = response.RatingSubmissonResponse{
		RatingID:     get.RatingID,
		UserID:       get.UserID,
		UserIDLegacy: get.UserIDLegacy,
		Comment:      get.Comment,
		Value:        get.Value,
	}
	return &result, message.SuccessMsg
}

// swagger:route GET /api/v1/rating-submissions RatingSubmission ListRatingSubmissionRequest
// Get List Rating Submissions
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) GetListRatingSubmissions(input request.ListRatingSubmissionRequest) ([]response.RatingSubmissonResponse, *base.Pagination, message.Message) {
	var dir interface{}
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
		input.Sort = "updated_at"
	}

	filter := request.RatingSubmissionFilter{}
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filter)
		if errMarshal != nil {
			return nil, nil, message.FailedMsg
		}
	}
	ratingSubmissions, pagination, err := s.ratingRepo.GetListRatingSubmissions(filter, input.Page, input.Limit, input.Sort, dir)
	if err != nil || ratingSubmissions == nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, message.ErrNoData
		}
	}

	results := make([]response.RatingSubmissonResponse, 0)
	for _, args := range ratingSubmissions {
		results = append(results, response.RatingSubmissonResponse{
			RatingID:     args.RatingID,
			UserID:       args.UserID,
			UserIDLegacy: args.UserIDLegacy,
			Comment:      args.Comment,
			Value:        args.Value,
		})
	}
	if len(ratingSubmissions) == 0 {
		return results, pagination, message.SuccessMsg
	}

	return results, pagination, message.SuccessMsg
}

// swagger:route POST /api/v1/rating-types-likert/ RatingTypesLikert createRatingTypeLikertRequest
// Create Likert Rating Types
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) CreateRatingTypeLikert(input request.SaveRatingTypeLikertRequest) message.Message {

	errMsg := validateNumStatement(input)
	if errMsg.Message != "" {
		return errMsg
	}
	check := true
	if input.Status == nil {
		input.Status = &check
	}

	err := s.ratingRepo.CreateRatingTypeLikert(input)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return message.ErrDuplicateType
		}
		return message.FailedMsg
	}
	return message.SuccessMsg
}

// swagger:route GET /api/v1/rating-types-likert/{id} RatingTypesLikert getRatingTypeLikertById
// Get Likert Rating Types by ID
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) GetRatingTypeLikertById(input request.GetRatingTypeLikertRequest) (*entity.RatingTypesLikertCol, message.Message) {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return nil, message.ErrNoData
	}
	result, err := s.ratingRepo.GetRatingTypeLikertById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.ErrNoData
		}
		return nil, message.FailedMsg
	}
	return result, message.SuccessMsg
}

// swagger:route PUT /api/v1/rating-types-likert/{id} RatingTypesLikert updateRatingTypeLikert
// Update Likert Rating Types
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) UpdateRatingTypeLikert(input request.SaveRatingTypeLikertRequest) message.Message {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return message.ErrNoData
	}

	_, err = s.ratingRepo.GetRatingTypeLikertById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return message.ErrNoData
		}
		return message.FailedMsg
	}

	rating, err := s.ratingRepo.GetRatingByType(input.Id)
	if err != nil && err != mongo.ErrNoDocuments {
		return message.FailedMsg
	}
	if rating != nil {
		errMsg := util.ValidInputUpdateRatingTypeLikertInRating(input)
		if errMsg.Message != "" {
			return errMsg
		}
		submissison, err := s.ratingRepo.GetRatingSubmissionByRatingId(rating.ID.Hex())
		if err != nil && err != mongo.ErrNoDocuments {
			return message.FailedMsg
		}
		if submissison != nil {
			errMsg = util.ValidInputUpdateRatingTypeLikertInSubmission(input)
			if errMsg.Message != "" {
				return errMsg
			}
		}
	}

	errMsg := validateNumStatement(input)
	if errMsg.Message != "" {
		return errMsg
	}

	err = s.ratingRepo.UpdateRatingTypeLikert(objectId, input)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return message.ErrDuplicateType
		}
		return message.FailedMsg
	}
	return message.SuccessMsg
}

// swagger:route DELETE /api/v1/rating-types-likert/{id} RatingTypesLikert deleteRatingTypeLikert
// Delete Likert Rating Types
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) DeleteRatingTypeLikertById(input request.GetRatingTypeLikertRequest) message.Message {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return message.ErrNoData
	}
	rating, err := s.ratingRepo.GetRatingByType(input.Id)
	if err != nil && err != mongo.ErrNoDocuments {
		return message.FailedMsg
	}
	if rating != nil {
		submissison, err := s.ratingRepo.GetRatingSubmissionByRatingId(rating.ID.Hex())
		if err != nil && err != mongo.ErrNoDocuments {
			return message.FailedMsg
		}
		if submissison != nil {
			return message.ErrThisRatingTypeIsInUse
		}
	}
	err = s.ratingRepo.DeleteRatingTypeLikert(objectId)
	if err != nil {
		if err.Error() == "user not found" {
			return message.ErrNoData
		}
		return message.FailedMsg
	}
	return message.SuccessMsg
}

// swagger:route GET /api/v1/rating-types-likert RatingTypesLikert getRatingTypeLikerts
// Get Likert Rating Types Listing
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) GetRatingTypeLikerts(input request.GetRatingTypeLikertsRequest) ([]entity.RatingTypesLikertCol, *base.Pagination, message.Message) {
	var dir interface{}
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
		input.Sort = "updated_at"
	}

	filter := request.FilterRatingTypeLikert{}
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filter)
		if errMarshal != nil {
			return nil, nil, message.FailedMsg
		}
	}

	ratingTypeLikerts, pagination, err := s.ratingRepo.GetRatingTypeLikerts(filter, input.Page, input.Limit, input.Sort, dir)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, message.ErrNoData
		}
		return nil, nil, message.FailedMsg
	}

	results := make([]entity.RatingTypesLikertCol, 0)
	if len(ratingTypeLikerts) == 0 {
		return results, pagination, message.SuccessMsg
	}
	results = ratingTypeLikerts

	return results, pagination, message.SuccessMsg
}

// swagger:route POST /api/v1/ratings/ Ratings CreateRatingRequest
// Create New Rating
//
// responses:
//  200: SuccessResponse
func (s *ratingServiceImpl) CreateRating(input request.SaveRatingRequest) (*entity.RatingsCol, message.Message) {

	if input.Status == nil {
		status := true
		input.Status = &status
	}

	// check source uid and source type not exist
	rating, err := s.ratingRepo.GetRatingBySourceUidAndSourceType(input.SourceUid, input.SourceType)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.FailedMsg
		}
	}

	if rating != nil {
		return nil, message.ErrExistingSourceUidAndSourceType
	}

	// check source exist
	source, err := s.medicalFacility.CallGetDetailMedicalFacility(input.SourceUid)
	if err != nil {
		return nil, message.ErrFailedToCallGetMedicalFacility
	}

	if source.Meta.Code != message.GetMedicalFacilitySuccess.Code {
		return nil, message.ErrSourceNotExist
	}

	// check rating type exist
	ratingTypeId, err := primitive.ObjectIDFromHex(input.RatingTypeId)
	if err != nil {
		return nil, message.ErrRatingTypeNotExist
	}

	ratingTypeNum, err := s.ratingRepo.GetRatingTypeNumByIdAndStatus(ratingTypeId)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.FailedMsg
		}
	}

	ratingTypeLikert, err := s.ratingRepo.GetRatingTypeLikertByIdAndStatus(ratingTypeId)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.FailedMsg
		}
	}

	if ratingTypeNum == nil && ratingTypeLikert == nil {
		return nil, message.ErrRatingTypeNotExist
	}

	if ratingTypeNum != nil && ratingTypeNum.Type != input.RatingType {
		return nil, message.ErrRatingTypeNotExist
	}

	if ratingTypeLikert != nil && ratingTypeLikert.Type != input.RatingType {
		return nil, message.ErrRatingTypeNotExist
	}

	result, err := s.ratingRepo.CreateRating(input)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, message.ErrDuplicateRatingName
		}
		return nil, message.FailedMsg
	}
	return result, message.SuccessMsg
}

// swagger:route GET /api/v1/ratings/{id} Ratings GetRatingRequest
// Get Rating By Id
//
// responses:
//  200: RatingsCol
func (s *ratingServiceImpl) GetRatingById(id string) (*entity.RatingsCol, message.Message) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, message.ErrDataNotFound
	}

	result, err := s.ratingRepo.GetRatingById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.ErrDataNotFound
		}
		return nil, message.FailedMsg
	}

	return result, message.SuccessMsg
}

// swagger:route PUT /api/v1/ratings/{id} Ratings UpdateRatingRequest
// Update Rating
//
// responses:
//  200: SuccessResponse
func (s *ratingServiceImpl) UpdateRating(input request.UpdateRatingRequest) message.Message {
	// get current rating
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return message.ErrDataNotFound
	}

	currentRating, err := s.ratingRepo.GetRatingById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return message.ErrDataNotFound
		}
		return message.FailedMsg
	}

	// check source uid and source type not exist
	if currentRating.SourceUid != input.Body.SourceUid || currentRating.SourceType != input.Body.SourceType {
		rating, err := s.ratingRepo.GetRatingBySourceUidAndSourceType(input.Body.SourceUid, input.Body.SourceType)
		if err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				return message.FailedMsg
			}
		}

		if rating != nil && rating.SourceUid == input.Body.SourceUid && rating.SourceType == input.Body.SourceType {
			return message.ErrExistingSourceUidAndSourceType
		}
	}

	// check source exist
	if input.Body.SourceUid != "" {
		source, err := s.medicalFacility.CallGetDetailMedicalFacility(input.Body.SourceUid)
		if err != nil {
			return message.ErrFailedToCallGetMedicalFacility
		}

		if source.Meta.Code != message.GetMedicalFacilitySuccess.Code {
			return message.ErrSourceNotExist
		}
	}

	_, err = s.ratingRepo.UpdateRating(objectId, input.Body)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return message.ErrDuplicateRatingName
		}
		return message.FailedMsg
	}
	return message.SuccessMsg
}

// swagger:route DELETE /api/v1/ratings/{id} Ratings DeleteRatingRequest
// Delete Rating
//
// responses:
//  200: SuccessResponse
func (s *ratingServiceImpl) DeleteRating(id string) message.Message {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return message.ErrDataNotFound
	}
	//FindRatingByRatingID

	rating, err := s.ratingRepo.GetRatingById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return message.ErrDataNotFound
		}
		return message.FailedMsg
	}

	ratingSubmission, err := s.ratingRepo.GetRatingSubmissionByRatingId(rating.ID.Hex())
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return message.FailedMsg
		}
	}

	if ratingSubmission != nil {
		return message.ErrRatingHasRatingSubmission
	}

	err = s.ratingRepo.DeleteRating(objectId)
	if err != nil {
		return message.FailedMsg
	}
	return message.SuccessMsg
}

// swagger:route GET /api/v1/ratings Ratings GetListRatingsRequest
// Get list Rating
//
// responses:
//  200: RatingsCol
func (s *ratingServiceImpl) GetListRatings(input request.GetListRatingsRequest) ([]entity.RatingsCol, *base.Pagination, message.Message) {
	input.MakeDefaultValueIfEmpty()
	var dir int
	if input.Dir == "asc" {
		dir = 1
	} else {
		dir = -1
	}

	filter := request.RatingFilter{}
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filter)
		if errMarshal != nil {
			return nil, nil, message.ErrUnmarshalFilterListRatingRequest
		}
	}

	ratings, pagination, err := s.ratingRepo.GetRatingsByParams(input.Limit, input.Page, dir, input.Sort, filter)
	if err != nil {
		return nil, nil, message.FailedMsg
	}
	results := make([]entity.RatingsCol, 0)
	if len(ratings) == 0 {
		return results, pagination, message.SuccessMsg
	}
	results = ratings

	return results, pagination, message.SuccessMsg
}

func validateNumStatement(input request.SaveRatingTypeLikertRequest) message.Message {
	var errMsg message.Message
	var myMap map[string]string
	count := 1
	datas, _ := json.Marshal(input)
	json.Unmarshal(datas, &myMap)

	for i := 1; i <= 10; i++ {
		var a string
		if i < 10 {
			a = fmt.Sprint("statement_0", i)
		} else {
			a = fmt.Sprint("statement_", i)
		}
		data := myMap[a]
		if data != "" {
			if count == input.NumStatements+1 {
				errMsg = message.ErrMatchNumState
				return errMsg
			}
			count++
		}
	}
	if count <= input.NumStatements {
		errMsg = message.ErrMatchNumState
		return errMsg
	}
	return errMsg
}

// swagger:route GET /api/v1/ratings/summary Ratings GetListRatingSummaryRequest
// Get list Rating Summary
//
// responses:
//  200: RatingsCol
func (s *ratingServiceImpl) GetListRatingSummary(input request.GetListRatingSummaryRequest) ([]response.RatingSummaryResponse, message.Message) {
	input.MakeDefaultValueIfEmpty()
	var dir int
	if input.Dir == "asc" {
		dir = 1
	} else {
		dir = -1
	}

	//find ratings by source_uid --> find rating submissions by id of rating
	filterForRating := request.RatingFilter{}
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filterForRating)
		if errMarshal != nil {
			return nil, message.ErrUnmarshalRequest
		}
	}

	if len(filterForRating.SourceUid) <= 0 {
		return nil, message.ErrSourceUidRequire
	}

	findR, _, errR := s.ratingRepo.GetRatingsByParams(input.Limit, input.Page, dir, input.Sort, filterForRating)
	if errR != nil {
		return nil, message.FailedMsg
	}

	//find rating summary, start_date(created_at), end_date(created_at), rating_id
	var rates = []string{}
	for _, args := range findR {
		rates = append(rates, args.ID.Hex())
	}

	var filterforRatingSub = request.RatingSubmissionFilter{
		RatingID: rates,
	}

	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filterforRatingSub)
		if errMarshal != nil {
			return nil, message.ErrUnmarshalRequest
		}
	}

	findS, _, err := s.ratingRepo.GetListRatingSubmissions(filterforRatingSub, input.Page, int64(input.Limit), input.Sort, dir)
	if err != nil {
		return nil, message.Message{
			Code:    message.ValidationFailCode,
			Message: "end_date can not before start_date",
		}
	}

	var min, max float64
	if filterforRatingSub.Score == nil {
		min, max = 0, 10
	} else {
		min, max = getScore(filterforRatingSub.Score)
	}

	// Handler response
	results := make([]response.RatingSummaryResponse, 0)
	for _, args := range filterForRating.SourceUid {
		var totalSub int
		var valueRate float64
		var avgValue float64
		for _, argRs := range findR {
			if argRs.SourceUid != args {
				continue
			}
			for _, argSs := range findS {
				if argSs.RatingID != argRs.ID.Hex() {
					continue
				}
				value, _ := strconv.ParseFloat(argSs.Value, 64)
				totalSub += 1
				valueRate += value
			}
			avgValue = valueRate / float64(totalSub)
			if math.IsNaN(avgValue) {
				avgValue = 0
			}
		}
		avgValue = util.RoundFloatWithPrecision(avgValue, 1)
		if filterScore(min, max, avgValue) {
			results = append(results, response.RatingSummaryResponse{
				SourceUID:   args,
				TotalReview: totalSub,
				Value:       avgValue,
			})
		}
	}

	if len(filterForRating.SourceUid) == 0 {
		return results, message.SuccessMsg
	}

	return results, message.SuccessMsg
}

func filterScore(min float64, max float64, value float64) bool {
	if value >= min && value <= max {
		return true
	}
	return false
}

func getScore(score []string) (float64, float64) {
	var min, max float64
	min, _ = strconv.ParseFloat(score[0], 64)
	min, _ = strconv.ParseFloat(score[1], 64)
	return min, max
}
