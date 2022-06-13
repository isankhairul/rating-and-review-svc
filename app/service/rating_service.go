package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/log"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
)

type RatingService interface {
	// Rating type num
	CreateRatingTypeNum(input request.CreateRatingTypeNumRequest) message.Message
	UpdateRatingTypeNum(input request.CreateRatingTypeNumRequest) message.Message
	GetRatingTypeNumById(input request.GetRatingTypeNumRequest) (*entity.RatingTypesNumCol, message.Message)
	DeleteRatingTypeNumById(input request.GetRatingTypeNumRequest) message.Message
	GetRatingTypeNums(input request.GetRatingTypeNumsRequest) ([]entity.RatingTypesNumCol, *base.Pagination, message.Message)

	// Rating submission
	CreateRatingSubmission(input request.CreateRatingSubmissonRequest) message.Message
	UpdateRatingSubmission(input request.UpdateRatingSubmissonRequest) message.Message
	GetRatingSubmission(id string) (*entity.RatingSubmisson, message.Message)
	GetListRatingSubmissions(input request.ListRatingSubmissionRequest) ([]entity.RatingSubmisson, *base.Pagination, message.Message)
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

// swagger:route POST /api/v1/rating-types-numeric ratingTypeNum createRatingTypeNum
// Create Rating Type Num
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) CreateRatingTypeNum(input request.CreateRatingTypeNumRequest) message.Message {
	if *input.Scale < 0 || *input.Scale > 2 {
		return message.ErrScaleValueReq
	}

	_, err := s.ratingRepo.CreateRatingTypeNum(input)
	if err != nil {
		return message.ErrSaveData
	}
	return message.SuccessMsg
}

// swagger:route GET /api/v1/rating-types-numeric/{id} ratingTypeNum getRatingById
// Get Rating Type Num By Id
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) GetRatingTypeNumById(input request.GetRatingTypeNumRequest) (*entity.RatingTypesNumCol, message.Message) {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return nil, message.ErrIdFormatReq
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

// swagger:route PUT /api/v1/rating-types-numeric/{id} ratingTypeNum updateRatingTypeNum
// Update Rating Type Num
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) UpdateRatingTypeNum(input request.CreateRatingTypeNumRequest) message.Message {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return message.ErrIdFormatReq
	}
	err = s.ratingRepo.UpdateRatingTypeNum(objectId, input)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return message.ErrNoData
		}
		return message.ErrSaveData
	}
	return message.SuccessMsg
}

// swagger:route Delete /api/v1/rating-types-numeric/{id} ratingTypeNum deleteRatingTypeNum
// Update Rating Type Num
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) DeleteRatingTypeNumById(input request.GetRatingTypeNumRequest) message.Message {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return message.ErrIdFormatReq
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

// swagger:route GET /api/v1/rating-types-numeric ratingTypeNum getRatingTypeNums
// Get Rating Type Nums
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
	result, pagination, err := s.ratingRepo.GetRatingTypeNums(filter, input.Page, input.Limit, sort, dir)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, message.ErrNoData
		}
		return nil, nil, message.FailedMsg
	}
	return result, pagination, message.SuccessMsg
}

// swagger:route POST /api/v1/rating-submissions ratingSubmission ReqRatingSubmissonBody
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
	if input.Value > 5.0 || input.Value < 0.0 {
		return message.ErrValueFormat
	}
	regexNumber := "[0-5]+"
	regexDecimal := "[+]?([0-5]*\\.[0]+|[5])"
	if findN.Intervals == 6 {
		matched, _ := regexp.MatchString(regexNumber, fmt.Sprintf("%f", input.Value))
		if matched == false {
			return message.ErrValueFormat
		}
	} else if findN.Intervals == 11 {
		matched, _ := regexp.MatchString(regexDecimal, fmt.Sprintf("%f", input.Value))
		if matched == false {
			return message.ErrValueFormat
		}
	}

	//The maximum length of user_agent allowed is 200 characters. Crop at 197 characters with triple dots (...) at the end.
	if len(input.UserAgent) > 200 {
		return message.UserAgentTooLong
	}

	// A submission with a combination of either (rating_id and user_id) OR (rating_id and user_id_legacy) is allowed once
	if input.UserID == nil {
		findL, errL := s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(input.UserIDLegacy, input.RatingID)
		if findL.UserIDLegacy != input.UserIDLegacy && errL == nil {
			return message.UserRated
		}
	}
	if input.UserIDLegacy == nil {
		find, errF := s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(input.UserID, input.RatingID)
		if find.UserID != input.UserID || errF == nil {
			return message.UserRated
		}
	}

	_, errC := s.ratingRepo.CreateRatingSubmission(input)
	if errC != nil {
		return message.ErrSaveData
	}
	return message.SuccessMsg

}

// swagger:route PUT /api/v1/rating-submissions/{id} ratingSubmission ReqUpdateRatingSubmissonBody
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
	if input.Value > 5.0 || input.Value < 0.0 {
		return message.ErrValueFormat
	}
	regexNumber := "[0-5]+"
	regexDecimal := "[+]?([0-5]*\\.[0]+|[5])"
	if findN.Intervals == 6 {
		matched, _ := regexp.MatchString(regexNumber, fmt.Sprintf("%f", input.Value))
		if matched == false {
			return message.ErrValueFormat
		}
	} else if findN.Intervals == 11 {
		matched, _ := regexp.MatchString(regexDecimal, fmt.Sprintf("%f", input.Value))
		if matched == false {
			return message.ErrValueFormat
		}
	}

	// A submission with a combination of either (rating_id and user_id) OR (rating_id and user_id_legacy) is allowed once
	if input.UserID == nil {
		objectUserLegacyId, _ := primitive.ObjectIDFromHex(*input.UserIDLegacy)
		findL, errL := s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(input.UserIDLegacy, input.RatingID)
		if (findL != nil && findL.ID != objectUserLegacyId) || errL == nil {
			return message.UserRated
		}
	}
	if input.UserIDLegacy == nil {
		objectUserId, _ := primitive.ObjectIDFromHex(*input.UserID)
		find, errF := s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(input.UserID, input.RatingID)
		if (find != nil && find.ID != objectUserId) || errF == nil {
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

// swagger:route DELETE /api/v1/rating-submissions/{id} ratingSubmission ReqDeleteRatingSubmissionById
// Delete Rating Submission
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) DeleteRatingSubmission(id string) message.Message {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return message.ErrIdFormatReq
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

// swagger:route GET /api/v1/rating-submissions/{id} ratingSubmission ReqRatingSubmissionById
// Get Rating Submission By Id
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) GetRatingSubmission(id string) (*entity.RatingSubmisson, message.Message) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, message.FailedMsg
	}
	result, err := s.ratingRepo.GetRatingSubmissionById(objectId)
	if err != nil {
		fmt.Println("err.Error()", err.Error())
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.ErrNoData
		}
		return nil, message.FailedMsg
	}
	return result, message.SuccessMsg
}

// swagger:route GET /api/v1/rating-submissions ratingSubmission ListRatingSubmissionRequest
// Get List Rating Submissions
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) GetListRatingSubmissions(input request.ListRatingSubmissionRequest) ([]entity.RatingSubmisson, *base.Pagination, message.Message) {
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
			return nil, nil, message.ErrUnmarshalRequest
		}
	}
	result, pagination, err := s.ratingRepo.GetListRatingSubmissions(filter, input.Page, input.Limit, input.Sort, dir)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, message.ErrNoData
		}
		return nil, nil, message.FailedMsg
	}
	return result, pagination, message.SuccessMsg
}

// swagger:route POST /api/v1/rating-types-likert/ RatingTypesLikert createRatingTypeLikertRequest
// Create Rating Type Likert
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) CreateRatingTypeLikert(input request.SaveRatingTypeLikertRequest) message.Message {
	err := s.ratingRepo.CreateRatingTypeLikert(input)
	if err != nil {
		return message.ErrSaveData
	}
	return message.SuccessMsg
}

// swagger:route GET /api/v1/rating-types-likert/{id} RatingTypesLikert getRatingTypeLikertById
// Get Rating Type Likert By Id
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) GetRatingTypeLikertById(input request.GetRatingTypeLikertRequest) (*entity.RatingTypesLikertCol, message.Message) {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return nil, message.ErrIdFormatReq
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
// Update Rating Type Likert
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) UpdateRatingTypeLikert(input request.SaveRatingTypeLikertRequest) message.Message {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return message.ErrIdFormatReq
	}
	err = s.ratingRepo.UpdateRatingTypeLikert(objectId, input)
	if err != nil {
		fmt.Println("err.Error()", err.Error())
		return message.ErrSaveData
	}
	return message.SuccessMsg
}

// swagger:route DELETE /api/v1/rating-types-likert/{id} RatingTypesLikert deleteRatingTypeLikert
// Delete Rating Type Likert
//
// security:
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingServiceImpl) DeleteRatingTypeLikertById(input request.GetRatingTypeLikertRequest) message.Message {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return message.ErrIdFormatReq
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
// Get Rating Type Likerts
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

	result, pagination, err := s.ratingRepo.GetRatingTypeLikerts(filter, input.Page, input.Limit, input.Sort, dir)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, message.ErrNoData
		}
		return nil, nil, message.FailedMsg
	}
	return result, pagination, message.SuccessMsg
}

// swagger:route POST /api/v1/ratings/ Ratings CreateRatingRequest
// Create New Rating
//
// responses:
//  200: SuccessResponse
func (s *ratingServiceImpl) CreateRating(input request.SaveRatingRequest) (*entity.RatingsCol, message.Message) {
	source, err := util.CallGetDetailMedicalFacility(input.SourceUid)
	if err != nil {
		return nil, message.ErrFailedToCallGetMedicalFacility
	}

	if source.Meta.Code != message.GetMedicalFacilitySuccess.Code {
		return nil, message.ErrSourceNotExist
	}

	ratingTypeId, err := primitive.ObjectIDFromHex(input.RatingTypeId)
	if err != nil {
		return nil, message.ErrRatingTypeIdFormatReq
	}

	rating, err := s.ratingRepo.GetRatingByName(input.Name)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.FailedMsg
		}
	}

	if rating != nil {
		return nil, message.ErrDuplicateRatingName
	}

	ratingType, err := s.ratingRepo.GetRatingTypeNumById(ratingTypeId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.ErrRatingTypeNotExist
		}
		return nil, message.FailedMsg
	}

	if ratingType == nil || ratingType.Type != input.RatingType {
		return nil, message.ErrRatingTypeNotExist
	}

	if !input.Status {
		input.Status = true
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
		return nil, message.ErrIdFormatReq
	}

	result, err := s.ratingRepo.GetRatingById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.ErrNoData
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
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return message.ErrIdFormatReq
	}

	_, err = s.ratingRepo.GetRatingById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return message.ErrNoData
		}
		return message.FailedMsg
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
		return message.ErrIdFormatReq
	}

	_, err = s.ratingRepo.GetRatingById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return message.ErrNoData
		}
		return message.FailedMsg
	}
	err = s.ratingRepo.DeleteRating(objectId)
	if err != nil {
		return message.FailedMsg
	}
	return message.SuccessMsg
}

// swagger:route GET /api/v1/ratings/ Ratings GetListRatingsRequest
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
			return nil, nil, message.ErrUnmarshalRequest
		}
	}

	result, pagination, err := s.ratingRepo.GetRatingsByParams(input.Limit, input.Page, dir, input.Sort, filter)
	if err != nil {
		return nil, nil, message.FailedMsg
	}

	return result, pagination, message.SuccessMsg
}
