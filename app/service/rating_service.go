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
	CreateRatingSubmission(input request.CreateRatingSubmissonRequest) message.Message
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
	GetListRatingSummary(input request.GetListRatingSummaryRequest) ([]response.RatingSummaryResponse, *base.Pagination, message.Message)
}

type ratingServiceImpl struct {
	logger     log.Logger
	ratingRepo repository.RatingRepository
}

func NewRatingService(
	lg log.Logger,
	rr repository.RatingRepository,
) RatingService {
	return &ratingServiceImpl{lg, rr}
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
	if *input.MaxScore < *input.MinScore {
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
			msg := util.ValidInputUpdateRatingTypeNum(input)
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
			return nil, nil, message.ErrUnmarshalFilterListRatingRequest
		}
	}
	ratingTypeNums, pagination, err := s.ratingRepo.GetRatingTypeNums(filter, input.Page, input.Limit, sort, dir)
	if err != nil {
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
func (s *ratingServiceImpl) CreateRatingSubmission(input request.CreateRatingSubmissonRequest) message.Message {

	// One of the following user_id and user_id_legacy must be filled.
	if input.UserID == nil && input.UserIDLegacy == nil {
		return message.UserUIDRequired
	}

	//The value must be valid according to requirements of rating type
	// Find rating_type_id by rating_id
	objecRatingtId, _ := primitive.ObjectIDFromHex(input.RatingID)
	findT, errT := s.ratingRepo.FindRatingByRatingID(objecRatingtId)
	if findT == nil || errT != nil {
		return message.ErrRatingNotFound
	}

	// Find intervals of Numeric Rating Type by rating_type_id
	objecRatingNumericTypeId, _ := primitive.ObjectIDFromHex(findT.RatingTypeId)
	findN, errN := s.ratingRepo.FindRatingNumericTypeByRatingTypeID(objecRatingNumericTypeId)
	if findN == nil || errN != nil {
		return message.ErrRatingNumericTypeNotFound
	}

	// The value must be valid according to requirements of rating type
	values := util.ValidValue(*findN.MinScore, *findN.MaxScore, *findN.Intervals, *findN.Scale)
	isInclude := util.ValidateValue(values, *input.Value)
	if isInclude == false {
		return message.Message{
			Code:    message.ValidationFailCode,
			Message: "value must be included in : " + fmt.Sprintf("%v", values),
		}
	}

	//The maximum length of user_agent allowed is 200 characters. Crop at 197 characters with triple dots (...) at the end.
	if len(strings.TrimSpace(input.UserAgent)) > 200 {
		return message.UserAgentTooLong
	}

	// A submission with a combination of either (rating_id and user_id) OR (rating_id and user_id_legacy) is allowed once
	var emptyStr = ""
	if input.UserID == nil || input.UserID == &emptyStr {
		findL, errL := s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(input.UserIDLegacy, input.RatingID)
		if findL != nil || errL == nil {
			return message.UserRated
		}
	}
	if input.UserIDLegacy == nil || input.UserIDLegacy == &emptyStr {
		find, errF := s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(input.UserID, input.RatingID)
		if find != nil && errF == nil {
			return message.UserRated
		}
	}

	_, errC := s.ratingRepo.CreateRatingSubmission(input)
	if errC != nil {
		return message.ErrSaveData
	}
	return message.SuccessMsg

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

	//The value must be valid according to requirements of rating type
	// Find rating_type_id by rating_id
	objecRatingtId, _ := primitive.ObjectIDFromHex(input.RatingID)
	findT, errT := s.ratingRepo.FindRatingByRatingID(objecRatingtId)
	if findT == nil || errT != nil {
		return message.ErrRatingNotFound
	}

	// Find intervals of Numeric Rating Type by rating_type_id
	objecRatingNumericTypeId, _ := primitive.ObjectIDFromHex(findT.RatingTypeId)
	findN, errN := s.ratingRepo.FindRatingNumericTypeByRatingTypeID(objecRatingNumericTypeId)
	if findN == nil || errN != nil {
		return message.ErrRatingNumericTypeNotFound
	}

	// The value must be valid according to requirements of rating type
	values := util.ValidValue(*findN.MinScore, *findN.MaxScore, *findN.Intervals, *findN.Scale)
	isInclude := util.ValidateValue(values, input.Value)
	if isInclude == false {
		return message.Message{
			Code:    message.ValidationFailCode,
			Message: "value must be included in : " + fmt.Sprintf("%v", values),
		}
	}

	// A submission with a combination of either (rating_id and user_id) OR (rating_id and user_id_legacy) is allowed once
	if input.UserID == nil {
		findL, errL := s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(input.UserIDLegacy, input.RatingID)
		if findL != nil || errL == nil {
			return message.UserRated
		}
	}
	if input.UserIDLegacy == nil {
		find, errF := s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(input.UserID, input.RatingID)
		if find != nil && errF == nil {
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
			return nil, nil, message.ErrUnmarshalFilterListRatingRequest
		}
	}
	ratingSubmissions, pagination, err := s.ratingRepo.GetListRatingSubmissions(filter, input.Page, input.Limit, input.Sort, dir)
	if err != nil || ratingSubmissions == nil {
		return nil, nil, message.FailedMsg
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
			return nil, nil, message.ErrUnmarshalFilterListRatingRequest
		}
	}

	ratingTypeLikerts, pagination, err := s.ratingRepo.GetRatingTypeLikerts(filter, input.Page, input.Limit, input.Sort, dir)
	if err != nil {
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
	// check source exist
	source, err := util.CallGetDetailMedicalFacility(input.SourceUid)
	if err != nil {
		return nil, message.ErrFailedToCallGetMedicalFacility
	}

	if source.Meta.Code != message.GetMedicalFacilitySuccess.Code {
		return nil, message.ErrSourceNotExist
	}

	// check duplicate name
	rating, err := s.ratingRepo.GetRatingByName(input.Name)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.FailedMsg
		}
	}

	if rating != nil {
		return nil, message.ErrDuplicateRatingName
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

	record, err := s.ratingRepo.GetRatingById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return message.ErrDataNotFound
		}
		return message.FailedMsg
	}

	// check source exist
	source, err := util.CallGetDetailMedicalFacility(input.Body.SourceUid)
	if err != nil {
		return message.ErrFailedToCallGetMedicalFacility
	}

	if source.Meta.Code != message.GetMedicalFacilitySuccess.Code {
		return message.ErrSourceNotExist
	}

	// check duplicate name if name changed
	if record.Name != input.Body.Name {
		rating, err := s.ratingRepo.GetRatingByName(input.Body.Name)
		if err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				return message.FailedMsg
			}
		}
		if rating != nil {
			return message.ErrDuplicateRatingName
		}
	}

	// check rating type exist
	ratingTypeId, err := primitive.ObjectIDFromHex(input.Body.RatingTypeId)
	if err != nil {
		return message.ErrRatingTypeNotExist
	}

	ratingTypeNum, err := s.ratingRepo.GetRatingTypeNumByIdAndStatus(ratingTypeId)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return message.FailedMsg
		}
	}

	ratingTypeLikert, err := s.ratingRepo.GetRatingTypeLikertByIdAndStatus(ratingTypeId)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return message.FailedMsg
		}
	}

	if ratingTypeNum == nil && ratingTypeLikert == nil {
		return message.ErrRatingTypeNotExist
	}

	if ratingTypeNum != nil && ratingTypeNum.Type != input.Body.RatingType {
		return message.ErrRatingTypeNotExist
	}

	if ratingTypeLikert != nil && ratingTypeLikert.Type != input.Body.RatingType {
		return message.ErrRatingTypeNotExist
	}

	_, err = s.ratingRepo.UpdateRating(objectId, input.Body)
	if err != nil {
		return message.ErrSaveData
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
	var data map[string]string
	request, _ := json.Marshal(input)
	json.Unmarshal(request, &data)

	for i := 1; i <= 10; i++ {
		var statement string
		if i < 10 {
			statement = fmt.Sprint("statement_0", i)
		} else {
			statement = fmt.Sprint("statement_", i)
		}
		dataStatement := data[statement]
		if i == input.NumStatements+1 {
			if dataStatement != "" {
				errMsg = message.ErrMatchNumState
				return errMsg
			}
		} else if i <= input.NumStatements {
			if dataStatement == "" {
				errMsg = message.ErrMatchNumState
				return errMsg
			}
		}
	}
	return errMsg
}

// swagger:route GET /api/v1/ratings/summary Ratings GetListRatingSummaryRequest
// Get list Rating Summary
//
// responses:
//  200: RatingsCol
func (s *ratingServiceImpl) GetListRatingSummary(input request.GetListRatingSummaryRequest) ([]response.RatingSummaryResponse, *base.Pagination, message.Message) {
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
			return nil, nil, message.ErrUnmarshalFilterListRatingRequest
		}
	}

	findRatings, _, errR := s.ratingRepo.GetRatingsByParams(input.Limit, input.Page, dir, input.Sort, filterForRating)
	if errR != nil {
		return nil, nil, message.FailedMsg
	}

	//find rating summary by score, start_date(created_at), end_date(created_at), rating_id
	var rates = []string{}
	for _, args := range findRatings {
		rates = append(rates, args.ID.Hex())
	}

	var filterForRatingSub = request.RatingSubmissionFilter{
		RatingID: rates,
	}

	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filterForRatingSub)
		if errMarshal != nil {
			return nil, nil, message.ErrUnmarshalFilterListRatingRequest
		}
	}

	findSubs, pagination, errS := s.ratingRepo.GetListRatingSubmissions(filterForRatingSub, input.Page, int64(input.Limit), input.Sort, dir)
	if errS != nil {
		return nil, nil, message.FailedMsg
	}

	results := make([]response.RatingSummaryResponse, 0)

	for _, args := range findSubs {
		for _, argsR := range findRatings {
			if argsR.ID.Hex() == args.RatingID {
				totalReviews := len(findSubs)
				results = append(results, response.RatingSummaryResponse{
					RatingID: args.RatingID,
					UserID:   args.UserID,
					//To do : get value
					UserIDLegacy: args.UserIDLegacy,
					TotalReview:  totalReviews,
					Value:        args.Value,
					SourceUID:    argsR.SourceUid,
					Name:         argsR.Name,
				})
			}
		}
	}
	if len(findSubs) == 0 {
		return results, pagination, message.SuccessMsg
	}

	return results, pagination, message.SuccessMsg
}
