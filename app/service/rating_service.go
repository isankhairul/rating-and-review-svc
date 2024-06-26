package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	publicrequest "go-klikdokter/app/model/request/public"
	"go-klikdokter/app/model/response"
	publicresponse "go-klikdokter/app/model/response/public"
	"go-klikdokter/app/repository"
	repository2 "go-klikdokter/app/repository/public"
	"go-klikdokter/helper/config"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RatingService interface {
	// Rating type num
	CreateRatingTypeNum(input request.CreateRatingTypeNumRequest) (*entity.RatingTypesNumCol, message.Message)
	UpdateRatingTypeNum(input request.EditRatingTypeNumRequest) message.Message
	GetRatingTypeNumById(input request.GetRatingTypeNumRequest) (*entity.RatingTypesNumCol, message.Message)
	DeleteRatingTypeNumById(input request.GetRatingTypeNumRequest) message.Message
	GetRatingTypeNums(input request.GetRatingTypeNumsRequest) ([]entity.RatingTypesNumCol, *base.Pagination, message.Message)

	// Rating submission
	CreateRatingSubmission(input request.CreateRatingSubmissionRequest) ([]publicresponse.PublicCreateRatingSubmissionResponse, message.Message)
	UpdateRatingSubmission(input request.UpdateRatingSubmissionRequest) message.Message
	GetRatingSubmission(id string) (*response.RatingSubmissonResponse, message.Message)
	GetListRatingSubmissions(input request.ListRatingSubmissionRequest) ([]response.RatingSubmissonResponse, *base.Pagination, message.Message)
	DeleteRatingSubmission(id string) message.Message
	GetListRatingSubmissionWithUserIdLegacy(input request.GetPublicListRatingSubmissionByUserIdRequest) ([]publicresponse.PublicRatingSubmissionResponse, *base.Pagination, message.Message)
	UpdateRatingSubDisplayNameByIdLegacy(input request.UpdateRatingSubDisplayNameRequest) message.Message
	CancelRatingSubmission(input request.CancelRatingById) message.Message

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
	GetRatingBySourceTypeAndActor(input publicrequest.GetRatingBySourceTypeAndActorRequest) (*publicresponse.RatingBySourceTypeAndActorResponse, message.Message)

	// Rating Formula
	CreateRatingFormula(input request.SaveRatingFormula) (*entity.RatingFormulaCol, message.Message)
	GetRatingFormulaById(input request.GetRatingFormulaRequest) (*entity.RatingFormulaCol, message.Message)
	UpdateRatingFormula(input request.SaveRatingFormula) message.Message
	DeleteRatingFormulaById(input request.GetRatingFormulaRequest) message.Message
	GetRatingFormulas(input request.GetRatingFormulasRequest) ([]entity.RatingFormulaCol, *base.Pagination, message.Message)

	// Rating Sub Helpful
	CreateRatingSubHelpful(input request.CreateRatingSubHelpfulRequest) (response.RatingSubHelpfulResponse, message.Message)
}

type ratingServiceImpl struct {
	logger           log.Logger
	ratingRepo       repository.RatingRepository
	publicRatingRepo repository2.PublicRatingRepository
	medicalFacility  util.MedicalFacilitySvc
}

func NewRatingService(
	lg log.Logger,
	rr repository.RatingRepository,
	prr repository2.PublicRatingRepository,
	mf util.MedicalFacilitySvc,
) RatingService {
	return &ratingServiceImpl{lg, rr, prr, mf}
}

// swagger:route POST /rating-types-numeric/ RatingTypeNum createRatingTypeNum
// Create Numeric Rating Types
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
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

// swagger:route GET /rating-types-numeric/{id} RatingTypeNum getRatingTypeNumById
// Get Numeric Rating Types by ID
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
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
		submission, err := s.ratingRepo.GetRatingSubmissionByRatingId(rating.ID.Hex())
		if err != nil && err != mongo.ErrNoDocuments {
			return message.FailedMsg
		}
		if submission != nil {
			msg := util.ValidInputUpdateRatingTypeNumSubmission(input)
			if msg != message.SuccessMsg {
				return msg
			}
		}
	}
	interval := util.ValidInterval(*input.MinScore, *input.MaxScore, *input.Scale)
	if *input.Intervals != interval {
		return message.Message{
			Code:    message.ValidationFailCode,
			Message: "interval must be " + strconv.Itoa(interval),
		}
	}

	return message.SuccessMsg
}

// swagger:route PUT /rating-types-numeric/{id} RatingTypeNum updateRatingTypeNum
// Update Numeric Rating Type
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
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

// swagger:route Delete /rating-types-numeric/{id} RatingTypeNum deleteRatingTypeNum
// Delete Numeric Rating Type
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
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
		submission, err := s.ratingRepo.GetRatingSubmissionByRatingId(rating.ID.Hex())
		if err != nil && err != mongo.ErrNoDocuments {
			return message.FailedMsg
		}
		if submission != nil {
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

// swagger:route GET /rating-types-numeric RatingTypeNum getRatingTypeNums
// Get Numeric Rating Types Listing
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
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

// swagger:route POST /rating-submissions/ RatingSubmission ReqRatingSubmissionBody
// Create Rating Submission
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) CreateRatingSubmission(input request.CreateRatingSubmissionRequest) ([]publicresponse.PublicCreateRatingSubmissionResponse, message.Message) {
	logger := log.With(s.logger, "RatingService", "CreateRatingSubmission")
	var saveReq = make([]request.SaveRatingSubmission, 0)
	var empty = ""
	var likertId = ""
	var valueLikert = ""
	var numId = ""
	result := []publicresponse.PublicCreateRatingSubmissionResponse{}
	haveRatNum := false
	haveRatLikert := false
	isAllRating := false
	valueLayanan := "1"
	isOrderIdExist := false

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

	// Validate displayname
	if input.DisplayName == nil || *input.DisplayName == "" {
		return result, message.ErrDisplayNameRequired
	}

	// The maximum length of user_agent allowed is 200 characters. Crop at 197 characters with triple dots (...) at the end.
	if len(strings.TrimSpace(input.UserAgent)) > 200 {
		return result, message.UserAgentTooLong
	}

	originalSourceTransID := input.SourceTransID
	if len(input.Ratings) == 0 && input.RatingType != "" {
		// Condition for All Rating (not doctor)
		isAllRating = true
		// Get Rating by SourceUID and RatingType
		rating, err := s.ratingRepo.FindRatingBySourceUIDAndRatingType(input.SourceUID, input.RatingType)
		if err != nil {
			return result, message.ErrRatingNotFound
		}
		if rating == nil || !*rating.Status {
			return result, message.ErrRatingNotFound
		}

		// Concate source_trans_id
		input.SourceTransID = originalSourceTransID + "||" + rating.ID.Hex()

		// A submission with a combination of either (rating_id and user_id) OR (rating_id and user_id_legacy) is allowed once
		userHasSubmitRating := checkUserHaveSubmitRating(*input.UserID, *input.UserIDLegacy, rating.ID.Hex(), input.SourceTransID, s)
		if userHasSubmitRating != message.SuccessMsg {
			return result, userHasSubmitRating
		}

		saveReq = append(saveReq, request.SaveRatingSubmission{
			RatingID:      rating.ID.Hex(),
			Value:         &valueLayanan,
			UserID:        input.UserID,
			UserIDLegacy:  input.UserIDLegacy,
			DisplayName:   input.DisplayName,
			Comment:       input.Comment,
			Avatar:        input.Avatar,
			IPAddress:     input.IPAddress,
			UserAgent:     input.UserAgent,
			SourceTransID: input.SourceTransID,
			UserPlatform:  input.UserPlatform,
			SourceUID:     input.SourceUID,
			IsAnonymous:   input.IsAnonymous,
			SourceType:    rating.SourceType,
		})
	} else {
		// Condition for Doctor Rating
		for _, argRatings := range input.Ratings {
			// Concate source_trans_id
			input.SourceTransID = originalSourceTransID + "||" + argRatings.ID

			// Find rating_type_id by rating_id
			objectRatingId, err := primitive.ObjectIDFromHex(argRatings.ID)
			if err != nil {
				return result, message.ErrRatingNotFound
			}
			rating, err := s.ratingRepo.FindRatingByRatingID(objectRatingId)
			if err != nil {
				return result, message.ErrRatingNotFound
			}
			if rating == nil || !*rating.Status {
				return result, message.ErrRatingNotFound
			}

			// Validate Numeric Type Value
			objectRatingTypeId, err := primitive.ObjectIDFromHex(rating.RatingTypeId)
			if err != nil {
				return result, message.ErrRatingNumericTypeNotFound
			}
			var validateMsg message.Message

			ratingNumericType, err := s.ratingRepo.FindRatingNumericTypeByRatingTypeID(objectRatingTypeId)
			if err != nil || !*ratingNumericType.Status {
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

			// Validate Likert Type Value
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
			userHasSubmitRating := checkUserHaveSubmitRating(*input.UserID, *input.UserIDLegacy, argRatings.ID, input.SourceTransID, s)
			if userHasSubmitRating != message.SuccessMsg {
				return result, userHasSubmitRating
			}

			if isExisted := isIdExisted(saveReq, argRatings.ID); !isExisted {
				return result, message.ErrCannotSameRatingId
			}

			if ratingNumericType == nil {
				haveRatLikert = true
				likertId = argRatings.ID
				valueLikert = *argRatings.Value
				saveReq = append(saveReq, request.SaveRatingSubmission{
					RatingID:      argRatings.ID,
					Value:         argRatings.Value,
					UserID:        input.UserID,
					UserIDLegacy:  input.UserIDLegacy,
					DisplayName:   input.DisplayName,
					Comment:       "",
					Avatar:        input.Avatar,
					IPAddress:     input.IPAddress,
					UserAgent:     input.UserAgent,
					SourceTransID: input.SourceTransID,
					UserPlatform:  input.UserPlatform,
					SourceUID:     input.SourceUID,
					IsAnonymous:   input.IsAnonymous,
				})
			} else {
				haveRatNum = true
				numId = argRatings.ID
				saveReq = append(saveReq, request.SaveRatingSubmission{
					RatingID:      argRatings.ID,
					Value:         argRatings.Value,
					UserID:        input.UserID,
					UserIDLegacy:  input.UserIDLegacy,
					DisplayName:   input.DisplayName,
					Comment:       input.Comment,
					Avatar:        input.Avatar,
					IPAddress:     input.IPAddress,
					UserAgent:     input.UserAgent,
					SourceTransID: input.SourceTransID,
					UserPlatform:  input.UserPlatform,
					SourceUID:     input.SourceUID,
					IsAnonymous:   input.IsAnonymous,
				})
			}
		}
	}

	// Set Tagging for Numeric Rating Sub
	if haveRatNum && haveRatLikert && !isAllRating {
		saveReq = createTagging(saveReq, numId, likertId, valueLikert)
	}

	// Check order_id exist for layanan
	if isAllRating {
		msg, err := util.CheckOrderIdExist(originalSourceTransID, logger)
		if err != nil {
			return result, message.ErrFailedRequestToPayment
		}

		if msg == message.SuccessMsg {
			isOrderIdExist = true
		} else {
			return result, msg
		}
	}

	if isOrderIdExist {
		go UpdateFlagPayment(originalSourceTransID, logger)
	}

	if len(saveReq) == 0 {
		return result, message.ErrTypeNotFound
	}
	ratingSubs, err := s.ratingRepo.CreateRatingSubmission(saveReq)
	if err != nil {
		return result, message.ErrSaveData
	}

	for _, arg := range *ratingSubs {
		data := publicresponse.PublicCreateRatingSubmissionResponse{}
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

func getSourceTypeByRatingType(ratingType string) string {
	sourceType := "doctor"
	if ratingType == "review_for_layanan" {
		sourceType = "layanan"
	}

	return sourceType
}

func UpdateFlagPayment(orderId string, logger log.Logger) error {
	msg, err := util.UpdateFlagPayment(orderId, logger)
	if err != nil {
		_ = level.Error(logger).Log("error", err)
		return err
	}

	if msg == message.SuccessMsg {
		return nil
	} else {
		_ = level.Error(logger).Log("error", "response code is not 200")
		return errors.New("response code request is not 200")
	}
}

func checkUserHaveSubmitRating(userId, userIdLegacy, ratingId, sourceTransId string, s *ratingServiceImpl) message.Message {
	if userId != "" && userIdLegacy != "" {
		ratingSubmission, er := s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(&userIdLegacy, ratingId, sourceTransId)
		val := util.ValidateUserIdAndUserIdLegacy(sourceTransId, ratingId, &userId, &userIdLegacy, ratingSubmission, er)
		if val {
			return message.UserRated
		}
		if !val {
			ratingSubmission, er = s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(&userId, ratingId, sourceTransId)
			valL := util.ValidateUserIdAndUserIdLegacy(sourceTransId, ratingId, &userId, &userIdLegacy, ratingSubmission, er)
			if valL {
				return message.UserRated
			}
		}
	}

	if userId == "" {
		ratingSubmission, er := s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(&userIdLegacy, ratingId, sourceTransId)
		if val := util.ValidateUserIdAndUserIdLegacy(sourceTransId, ratingId, &userId, &userIdLegacy, ratingSubmission, er); val {
			return message.UserRated
		}
	}

	if userIdLegacy == "" {
		ratingSubmission, er := s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(&userId, ratingId, sourceTransId)
		if val := util.ValidateUserIdAndUserIdLegacy(sourceTransId, ratingId, &userId, &userIdLegacy, ratingSubmission, er); val {
			return message.UserRated
		}
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

func createTagging(saveReq []request.SaveRatingSubmission, numId, likertId, valueLikert string) []request.SaveRatingSubmission {
	// tagging := request.TaggingObj{
	// 	RatingId: likertId,
	// 	Value:    strings.Split(valueLikert, ","),
	// }
	for i, args := range saveReq {
		if numId == args.RatingID {
			saveReq[i].Tagging.RatingId = likertId
			saveReq[i].Tagging.Value = strings.Split(valueLikert, ",")
		}
	}
	return saveReq
}

// swagger:route PUT /rating-submissions/{id} RatingSubmission ReqUpdateRatingSubmissionBody
// Update Rating Submission
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) UpdateRatingSubmission(input request.UpdateRatingSubmissionRequest) message.Message {
	// Input ID of Submission
	objectRatingSubmissionId, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return message.RatingSubmissionNotFound
	}

	// Input Rating ID of rating ID
	objectRatingId, err := primitive.ObjectIDFromHex(input.RatingID)
	if err != nil {
		return message.ErrRatingNotFound
	}

	empty := ""

	// find Rating submission
	ratingSubmission, err := s.ratingRepo.GetRatingSubmissionById(objectRatingSubmissionId)
	if err != nil || ratingSubmission == nil {
		return message.RatingSubmissionNotFound
	}

	// Find rating_type_id by rating
	rating, err := s.ratingRepo.FindRatingByRatingID(objectRatingId)
	if err != nil {
		return message.ErrRatingNotFound
	}
	if rating == nil || *rating.Status == false {
		return message.ErrRatingNotFound
	}

	// Validate value of numeric type
	objectRatingTypeId, err := primitive.ObjectIDFromHex(rating.RatingTypeId)
	if err != nil {
		return message.ErrTypeNotFound
	}
	var validateMsg message.Message

	ratingNumericType, err := s.ratingRepo.FindRatingNumericTypeByRatingTypeID(objectRatingTypeId)
	if err != nil || *ratingNumericType.Status == false {
		validateMsg = message.ErrRatingNumericTypeNotFound
	} else {
		value, er := strconv.ParseFloat(*input.Value, 64)
		if er != nil {
			return message.ErrValueFormatForNumericType
		}
		validateMsg = util.ValidateTypeNumeric(ratingNumericType, value)
	}

	if validateMsg.Code == message.ValidationFailCode {
		return validateMsg
	}

	// Validate value of likert type
	if validateMsg == message.ErrRatingNumericTypeNotFound {
		ratingTypeLikert, er := s.ratingRepo.GetRatingTypeLikertByIdAndStatus(objectRatingTypeId)
		if er != nil {
			if errors.Is(mongo.ErrNoDocuments, er) {
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
		input.Comment = empty
	}

	// A submission with a combination of either (rating_id and user_id) OR (rating_id and user_id_legacy) is allowed once
	if ratingSubmission.UserID == nil {
		ratingSubmission.UserID = &empty
	}

	if ratingSubmission.UserIDLegacy == nil {
		ratingSubmission.UserIDLegacy = &empty
	}

	if ratingSubmission.UserID != &empty && ratingSubmission.UserIDLegacy != &empty {
		ratingSubmissionV, err := s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(ratingSubmission.UserIDLegacy, input.RatingID, ratingSubmission.SourceTransID)
		val := util.ValidateUserIdAndUserIdLegacyForUpdate(input, objectRatingSubmissionId, ratingSubmission.SourceTransID, ratingSubmission.UserID, ratingSubmission.UserIDLegacy, ratingSubmissionV, err)
		if val {
			return message.UserRated
		}
		if !val {
			ratingSubmissionV2, err := s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(ratingSubmission.UserID, input.RatingID, ratingSubmission.SourceTransID)
			val = util.ValidateUserIdAndUserIdLegacyForUpdate(input, objectRatingSubmissionId, ratingSubmission.SourceTransID, ratingSubmission.UserID, ratingSubmission.UserIDLegacy, ratingSubmissionV2, err)
			if val {
				return message.UserRated
			}
		}
	}

	if ratingSubmission.UserID == &empty {
		ratingSubmissionV, err := s.ratingRepo.FindRatingSubmissionByUserIDLegacyAndRatingID(ratingSubmission.UserIDLegacy, ratingSubmission.RatingID, ratingSubmission.SourceTransID)
		val := util.ValidateUserIdAndUserIdLegacyForUpdate(input, objectRatingSubmissionId, ratingSubmission.SourceTransID, ratingSubmission.UserID, ratingSubmission.UserIDLegacy, ratingSubmissionV, err)
		if val {
			return message.UserRated
		}
	}

	if ratingSubmission.UserIDLegacy == &empty {
		ratingSubmissionV, err := s.ratingRepo.FindRatingSubmissionByUserIDAndRatingID(ratingSubmission.UserID, input.RatingID, ratingSubmission.SourceTransID)
		val := util.ValidateUserIdAndUserIdLegacyForUpdate(input, objectRatingSubmissionId, ratingSubmission.SourceTransID, ratingSubmission.UserID, ratingSubmission.UserIDLegacy, ratingSubmissionV, err)
		if val {
			return message.UserRated
		}
	}

	// Update
	errC := s.ratingRepo.UpdateRatingSubmission(input, objectRatingSubmissionId)
	if errC != nil {
		return message.ErrSaveData
	}
	return message.SuccessMsg

}

// swagger:route DELETE /rating-submissions/{id} RatingSubmission ReqDeleteRatingSubmissionById
// Delete Rating Submission
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
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

// swagger:route GET /list-rating-submissions/{source_type}/{source_uid}/{user_id_legacy} RatingSubmission GetPublicListRatingSubmissionByUserIdRequest
// Get Rating Submission with User Id Legacy
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) GetListRatingSubmissionWithUserIdLegacy(input request.GetPublicListRatingSubmissionByUserIdRequest) ([]publicresponse.PublicRatingSubmissionResponse, *base.Pagination, message.Message) {
	results := []publicresponse.PublicRatingSubmissionResponse{}
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
				Code:    message.DataNotFoundCode,
				Message: "Cannot find rating with params SourceType :" + input.SourceType + ", SourceUid:" + input.SourceUID,
			}
		}
		return nil, nil, message.FailedMsg
	}
	if len(ratings) <= 0 {
		return nil, nil, message.Message{
			Code:    message.DataNotFoundCode,
			Message: "Cannot find rating with params SourceType :" + input.SourceType + ", SourceUid:" + input.SourceUID,
		}
	}

	// Get Rating Submission
	filterRatingSubs := publicrequest.FilterRatingSubmission{}
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filterRatingSubs)
		if errMarshal != nil {
			return nil, nil, message.ErrUnmarshalFilterListRatingRequest
		}
	}

	for _, v := range ratings {
		filterRatingSubs.RatingID = append(filterRatingSubs.RatingID, v.ID.Hex())
	}

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
		likeByme := false
		// Get Rating value
		ratingId, err := primitive.ObjectIDFromHex(v.RatingID)
		if err != nil {
			return nil, nil, message.FailedMsg
		}
		rating, err := s.ratingRepo.GetRatingById(ratingId)
		if err != nil {
			return nil, nil, message.ErrRatingNotFound
		}

		// Get like by me value
		ratingSubHelpful, err := s.publicRatingRepo.GetRatingSubHelpfulByRatingSubAndActor(v.ID.Hex(), input.UserIdLegacy)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return results, pagination, message.ErrNoData
			}
			return nil, nil, message.FailedMsg
		}
		if ratingSubHelpful != nil {
			likeByme = ratingSubHelpful.Status
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
		results = append(results, publicresponse.PublicRatingSubmissionResponse{
			ID:            v.ID,
			UserID:        v.UserID,
			UserIDLegacy:  v.UserIDLegacy,
			DisplayName:   displayName,
			Avatar:        v.Avatar,
			Comment:       v.Comment,
			SourceTransID: v.SourceTransID,
			LikeCounter:   v.LikeCounter,
			LikeByMe:      likeByme,
			RatingType:    rating.RatingType,
			Value:         v.Value,
			CreatedAt:     v.CreatedAt.In(Loc),
		})
	}
	return results, pagination, message.SuccessMsg
}

// swagger:route PUT /rating-submissions/user-id-legacy/{user_id_legacy} RatingSubmission ReqUpdateRatingSubDisplayNameBody
// Update Rating Submission Display Name By Id Legacy
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) UpdateRatingSubDisplayNameByIdLegacy(input request.UpdateRatingSubDisplayNameRequest) message.Message {
	if input.DisplayName == "" {
		return message.ErrDisplayNameRequired
	}

	err := s.publicRatingRepo.UpdateRatingSubDisplayNameByIdLegacy(input)
	if err != nil {
		return message.ErrSaveData
	}

	return message.SuccessMsg
}

// swagger:route PUT /cancel/rating-submissions RatingSubmission ReqCancelRatingSubmission
// Cancel Rating Submission By Submission Id
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) CancelRatingSubmission(input request.CancelRatingById) message.Message {
	ids := []primitive.ObjectID{}
	reason := input.CancelledReason

	for _, id := range input.RatingSubmissionId {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return message.ErrRatingSubmissionNotFound
		}
		ids = append(ids, objectId)
	}

	err := s.ratingRepo.CancelRatingSubmissionByIds(ids, reason)
	if err != nil {
		return message.ErrSaveData
	}
	return message.SuccessMsg
}

// swagger:route GET /rating-submissions/{id} RatingSubmission ReqRatingSubmissionById
// Get Rating Submission By Id
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) GetRatingSubmission(id string) (*response.RatingSubmissonResponse, message.Message) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, message.ErrRatingSubmissionNotFound
	}
	get, err := s.ratingRepo.GetRatingSubmissionById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.ErrRatingSubmissionNotFound
		}
		return nil, message.FailedMsg
	}
	var result = response.RatingSubmissonResponse{
		RatingID:     get.RatingID,
		UserID:       get.UserID,
		UserIDLegacy: get.UserIDLegacy,
		Comment:      *get.Comment,
		Value:        get.Value,
		SourTransID:  get.SourceTransID,
	}
	return &result, message.SuccessMsg
}

// swagger:route GET /rating-submissions RatingSubmission ListRatingSubmissionRequest
// Get List Rating Submissions
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) GetListRatingSubmissions(input request.ListRatingSubmissionRequest) ([]response.RatingSubmissonResponse, *base.Pagination, message.Message) {
	var dir interface{}
	results := make([]response.RatingSubmissonResponse, 0)
	userIdEmpty := ""
	commentEmpty := ""
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

	filter := request.RatingSubmissionFilter{}
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filter)
		if errMarshal != nil {
			return nil, nil, message.WrongFilter
		}
	}
	ratingSubmissions, pagination, err := s.ratingRepo.GetListRatingSubmissions(filter, input.Page, input.Limit, input.Sort, dir)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return results, pagination, message.ErrDataNotFound
		}
		return nil, nil, message.FailedMsg
	}

	for _, args := range ratingSubmissions {
		if args.UserID == nil {
			args.UserID = &userIdEmpty
		}
		if args.Comment == nil {
			args.Comment = &commentEmpty
		}
		if filScore := filterScoreSubmission(args, filter.Score); filScore {
			results = append(results, response.RatingSubmissonResponse{
				RatingID:     args.RatingID,
				UserID:       args.UserID,
				UserIDLegacy: args.UserIDLegacy,
				Comment:      *args.Comment,
				Value:        args.Value,
				SourTransID:  args.SourceTransID,
			})
		}
	}
	if len(filter.Score) > 0 && pagination != nil {
		pagination.Records = int64(len(results))
		pagination.TotalRecords = int64(len(results))
	}

	return results, pagination, message.SuccessMsg
}

func filterScoreSubmission(ratingSubmissions entity.RatingSubmisson, score []float64) bool {
	if len(score) == 0 {
		return true
	}
	var scoreDB = make([]float64, 0)
	values := strings.Split(ratingSubmissions.Value, ",")
	for _, argVs := range values {
		value, _ := strconv.ParseFloat(argVs, 64)
		scoreDB = append(scoreDB, value)
	}
	for _, argsF := range score {
		for _, argSs := range scoreDB {
			if argSs == argsF {
				return true
			}
		}
	}
	return false
}

// swagger:route POST /rating-types-likert/ RatingTypesLikert createRatingTypeLikertRequest
// Create Likert Rating Types
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
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

// swagger:route GET /rating-types-likert/{id} RatingTypesLikert getRatingTypeLikertById
// Get Likert Rating Types by ID
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
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

// swagger:route PUT /rating-types-likert/{id} RatingTypesLikert updateRatingTypeLikert
// Update Likert Rating Types
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) UpdateRatingTypeLikert(input request.SaveRatingTypeLikertRequest) message.Message {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return message.ErrNoData
	}

	ratingTypeLikert, err := s.ratingRepo.GetRatingTypeLikertById(objectId)
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
		submission, err := s.ratingRepo.GetRatingSubmissionByRatingId(rating.ID.Hex())
		if err != nil && err != mongo.ErrNoDocuments {
			return message.FailedMsg
		}
		if submission == nil {
			errMsg := updateRatingTypeLikertHaveRating(s, input, ratingTypeLikert, objectId)
			if errMsg != message.SuccessMsg {
				return errMsg
			}
			return message.SuccessMsg
		}

		errMsg := updateRatingTypeLikertHaveSubmission(s, input, ratingTypeLikert, objectId)
		if errMsg != message.SuccessMsg {
			return errMsg
		}
		return message.SuccessMsg
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

// swagger:route DELETE /rating-types-likert/{id} RatingTypesLikert deleteRatingTypeLikert
// Delete Likert Rating Types
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
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
		submission, err := s.ratingRepo.GetRatingSubmissionByRatingId(rating.ID.Hex())
		if err != nil && err != mongo.ErrNoDocuments {
			return message.FailedMsg
		}
		if submission != nil {
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

// swagger:route GET /rating-types-likert RatingTypesLikert getRatingTypeLikerts
// Get Likert Rating Types Listing
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
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

// swagger:route POST /ratings/ Ratings CreateRatingRequest
// Create New Rating
//
// security:
// - Bearer: []
// responses:
//
//	200: SuccessResponse
func (s *ratingServiceImpl) CreateRating(input request.SaveRatingRequest) (*entity.RatingsCol, message.Message) {

	if input.Status == nil {
		status := true
		input.Status = &status
	}

	// check rating type id, source uid and source type not exist
	rating, err := s.ratingRepo.GetRatingByRatingTypeSourceUidAndSourceType(input.RatingTypeId, input.SourceUid, input.SourceType)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.FailedMsg
		}
	}
	if rating != nil {
		return nil, message.ErrExistingRatingTypeIdSourceUidAndSourceType
	}

	// Check Source UID for Hospital
	if input.SourceType == "hospital" {
		// check source exist
		source, err := s.medicalFacility.CallGetDetailMedicalFacility(input.SourceUid)
		if err != nil {
			return nil, message.ErrFailedToCallGetMedicalFacility
		}

		if source.Meta.Code != message.GetMedicalFacilitySuccess.Code {
			return nil, message.ErrSourceNotExist
		}
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

// swagger:route GET /ratings/{id} Ratings GetRatingRequest
// Get Rating By Id
//
// security:
// - Bearer: []
// responses:
//
//	200: SuccessResponse
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

// swagger:route PUT /ratings/{id} Ratings UpdateRatingRequest
// Update Rating
//
// security:
// - Bearer: []
// responses:
//
//	200: SuccessResponse
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

	// check rating has rating submission
	ratingSubmission, err := s.ratingRepo.GetRatingSubmissionByRatingId(input.Id)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return message.FailedMsg
		}
	}

	if ratingSubmission != nil && (input.Body.SourceUid != "" || input.Body.SourceType != "") {
		return message.ErrCanNotUpdateSourceTypeOrSoureUid
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

	// check source uid and source type not exist
	if (input.Body.SourceUid != "" && currentRating.SourceUid != input.Body.SourceUid) || (input.Body.SourceType != "" && currentRating.SourceType != input.Body.SourceType) {
		var sourceUid = input.Body.SourceUid
		var sourceType = input.Body.SourceType
		if sourceUid == "" {
			sourceUid = currentRating.SourceUid
		}
		if sourceType == "" {
			sourceType = currentRating.SourceType
		}
		rating, err := s.ratingRepo.GetRatingByRatingTypeSourceUidAndSourceType(currentRating.RatingTypeId, sourceUid, sourceType)
		if err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				return message.FailedMsg
			}
		}

		if rating != nil {
			return message.ErrExistingRatingTypeIdSourceUidAndSourceType
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

// swagger:route DELETE /ratings/{id} Ratings DeleteRatingRequest
// Delete Rating
//
// security:
// - Bearer: []
// responses:
//
//	200: SuccessResponse
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

// swagger:route GET /ratings Ratings GetListRatingsRequest
// Get list Rating
//
// security:
// - Bearer: []
// responses:
//
//	200: SuccessResponse
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
		var statement string
		if i < 10 {
			statement = fmt.Sprint("statement_0", i)
		} else {
			statement = fmt.Sprint("statement_", i)
		}
		data := myMap[statement]
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

func updateRatingTypeLikertHaveRating(s *ratingServiceImpl, input request.SaveRatingTypeLikertRequest, ratingTypeLikert *entity.RatingTypesLikertCol, objectId primitive.ObjectID) message.Message {
	errMsg := util.ValidInputUpdateRatingTypeLikertInRating(input)
	if errMsg.Message != "" {
		return errMsg
	}
	input.Type = ratingTypeLikert.Type
	input.Status = ratingTypeLikert.Status
	errMsg = validateNumStatement(input)
	if errMsg.Message != "" {
		return errMsg
	}
	err := s.ratingRepo.UpdateRatingTypeLikert(objectId, input)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return message.ErrDuplicateType
		}
		return message.FailedMsg
	}
	return message.SuccessMsg
}

func updateRatingTypeLikertHaveSubmission(s *ratingServiceImpl, input request.SaveRatingTypeLikertRequest, ratingTypeLikert *entity.RatingTypesLikertCol, objectId primitive.ObjectID) message.Message {
	errMsg := util.ValidInputUpdateRatingTypeLikertInSubmission(input)
	if errMsg.Message != "" {
		return errMsg
	}

	input.Type = ratingTypeLikert.Type
	input.NumStatements = ratingTypeLikert.NumStatements
	input.Statement01 = ratingTypeLikert.Statement01
	input.Statement02 = ratingTypeLikert.Statement02
	input.Statement03 = ratingTypeLikert.Statement03
	input.Statement04 = ratingTypeLikert.Statement04
	input.Statement05 = ratingTypeLikert.Statement05
	input.Statement06 = ratingTypeLikert.Statement06
	input.Statement07 = ratingTypeLikert.Statement07
	input.Statement08 = ratingTypeLikert.Statement08
	input.Statement09 = ratingTypeLikert.Statement09
	input.Statement10 = ratingTypeLikert.Statement10
	input.Status = ratingTypeLikert.Status
	errMsg = validateNumStatement(input)
	if errMsg.Message != "" {
		return errMsg
	}
	err := s.ratingRepo.UpdateRatingTypeLikert(objectId, input)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return message.ErrDuplicateType
		}
		return message.FailedMsg
	}
	return message.SuccessMsg
}

// swagger:route GET /ratings/summary/{source_type} Ratings GetListRatingSummaryRequest
// Get list Rating Summary
//
// security:
// - Bearer: []
// responses:
//
//	200: SuccessResponse
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
	filterForRating.SourceType = input.SourceType
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
	var rates = []string(nil)
	for _, args := range findR {
		rates = append(rates, args.ID.Hex())
	}

	var filterForRatingSub = request.RatingSubmissionFilter{
		RatingID: rates,
	}

	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filterForRatingSub)
		if errMarshal != nil {
			return nil, message.ErrUnmarshalRequest
		}
	}

	findS, _, err := s.ratingRepo.GetListRatingSubmissions(filterForRatingSub, input.Page, int64(input.Limit), input.Sort, dir)
	if err != nil {
		return nil, message.Message{
			Code:    message.ValidationFailCode,
			Message: "Wrong filter",
		}
	}

	var min, max float64
	if filterForRatingSub.Score == nil {
		min, max = 0, 10
	} else {
		if len(filterForRatingSub.Score) > 2 || len(filterForRatingSub.Score) < 2 {
			return nil, message.WrongScoreFilter
		}
		min, max = getScore(filterForRatingSub.Score)
	}

	results := make([]response.RatingSummaryResponse, 0)

	if len(filterForRating.SourceUid) == 0 || len(findR) == 0 || len(findS) == 0 {
		return results, message.SuccessMsg
	}
	// Handler response
	results = CalculateValue(filterForRating, findR, findS, min, max)

	return results, message.SuccessMsg
}

func CalculateValue(filterForRating request.RatingFilter, rating []entity.RatingsCol, ratingSubmission []entity.RatingSubmisson, min float64, max float64) []response.RatingSummaryResponse {
	results := make([]response.RatingSummaryResponse, 0)
	for _, args := range filterForRating.SourceUid {
		var totalSub int
		var valueRate float64
		var avgValue float64 = 0
		var totalReview int
		for _, argRs := range rating {
			if argRs.SourceUid != args {
				continue
			}
			for _, argSs := range ratingSubmission {
				if argSs.RatingID != argRs.ID.Hex() {
					continue
				}
				var avgValuePerSub float64
				var total float64
				values := make([]float64, 0)
				strValue := strings.Split(argSs.Value, ",")
				for _, argVs := range strValue {
					value, _ := strconv.ParseFloat(argVs, 64)
					values = append(values, value)
				}
				for _, argVPS := range values {
					total += argVPS
				}
				avgValuePerSub = total / float64(len(values))
				totalSub += len(values)
				valueRate += avgValuePerSub
			}
			totalReview = len(ratingSubmission)
			if !math.IsNaN(avgValue) {
				avgValue = valueRate / float64(len(ratingSubmission))
			}
		}
		avgValue = util.RoundFloatWithPrecision(avgValue, 1)
		if filterScore(min, max, avgValue) {
			results = append(results, response.RatingSummaryResponse{
				SourceUID:   args,
				TotalReview: totalReview,
				Value:       avgValue,
			})
		}
	}
	return results
}

func filterScore(min float64, max float64, value float64) bool {
	if value >= min && value <= max {
		return true
	}
	return false
}

func getScore(score []float64) (float64, float64) {
	var min, max float64
	min = score[0]
	max = score[1]
	return min, max
}

// swagger:route GET /list-ratings/{source_type}/{source_uid} Ratings GetRatingBySourceTypeAndActor
// Get Rating By Source Type and Source UID
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) GetRatingBySourceTypeAndActor(input publicrequest.GetRatingBySourceTypeAndActorRequest) (*publicresponse.RatingBySourceTypeAndActorResponse, message.Message) {
	result := publicresponse.RatingBySourceTypeAndActorResponse{}
	filter := publicrequest.GetRatingBySourceTypeAndActorFilter{}

	_ = json.Unmarshal([]byte(input.Filter), &filter)

	// Get Ratings By Type and Actor UID
	ratings, err := s.publicRatingRepo.GetRatingsBySourceTypeAndActor(input.SourceType, input.SourceUID, filter)
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
			likertResp := publicresponse.MapRatingLikertToRatingNumericResp(*likert, v.ID.Hex())
			result.Ratings = append(result.Ratings, likertResp)
		} else {
			numeric, err := s.publicRatingRepo.GetRatingTypeNumById(ratingTypeId)
			if err != nil {
				return nil, message.FailedMsg
			}
			if numeric == nil {
				return nil, message.ErrRatingTypeNotExist
			}
			numericResp := publicresponse.MapRatingNumericToRatingNumericResp(*numeric, v.ID.Hex())
			result.Ratings = append(result.Ratings, numericResp)
		}
	}
	return &result, message.SuccessMsg
}

// swagger:route POST /rating-formula/ RatingFormula createRatingFormula
// Create Rating Formula
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) CreateRatingFormula(input request.SaveRatingFormula) (*entity.RatingFormulaCol, message.Message) {
	check := true
	if input.Status == nil {
		input.Status = &check
	}
	result, err := s.ratingRepo.CreateRatingFormula(input)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, message.ErrDuplicateType
		}
		return nil, message.FailedMsg
	}
	return result, message.SuccessMsg
}

// swagger:route GET /rating-formula/{id} RatingFormula getRatingFormulaById
// Get Rating Formula by ID
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) GetRatingFormulaById(input request.GetRatingFormulaRequest) (*entity.RatingFormulaCol, message.Message) {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return nil, message.ErrNoData
	}
	result, err := s.ratingRepo.GetRatingFormulaById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.ErrNoData
		}
		return nil, message.FailedMsg
	}
	return result, message.SuccessMsg
}

// swagger:route PUT /rating-formula/{id} RatingFormula updateRatingFormula
// Update Rating Formula
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) UpdateRatingFormula(input request.SaveRatingFormula) message.Message {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return message.ErrNoData
	}
	err = s.ratingRepo.UpdateRatingFormula(objectId, input)
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

// swagger:route Delete /rating-formula/{id} RatingFormula deleteRatingFormula
// Delete Rating Formula
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) DeleteRatingFormulaById(input request.GetRatingFormulaRequest) message.Message {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return message.ErrNoData
	}
	err = s.ratingRepo.DeleteRatingFormula(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return message.ErrNoData
		}
		return message.FailedMsg
	}
	return message.SuccessMsg
}

// swagger:route GET /rating-formula RatingFormula getRatingFormulas
// Get Rating Formula Listing
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) GetRatingFormulas(input request.GetRatingFormulasRequest) ([]entity.RatingFormulaCol, *base.Pagination, message.Message) {
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

	filter := request.RatingFormulaFilter{}
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filter)
		if errMarshal != nil {
			return nil, nil, message.ErrUnmarshalRequest
		}
	}
	ratingFormulas, pagination, err := s.ratingRepo.GetRatingFormulas(filter, input.Page, input.Limit, sort, dir)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, message.ErrNoData
		}
		return nil, nil, message.FailedMsg
	}
	results := make([]entity.RatingFormulaCol, 0)
	if len(ratingFormulas) == 0 {
		return results, pagination, message.SuccessMsg
	}
	results = ratingFormulas

	return results, pagination, message.SuccessMsg
}

// swagger:route POST /helpful-rating-submission/ RatingSubHelpful ReqRatingSubHelpfulBody
// Create Helpful Rating Submission
//
// security:
// - Bearer: []
// responses:
//
//	401: SuccessResponse
//	200: SuccessResponse
func (s *ratingServiceImpl) CreateRatingSubHelpful(input request.CreateRatingSubHelpfulRequest) (response.RatingSubHelpfulResponse, message.Message) {
	result := response.RatingSubHelpfulResponse{}
	result.RatingSubmissionId = input.RatingSubmissionID
	result.UserIdLegacy = input.UserIDLegacy
	counter := 0
	// check rating submission is exist
	ratingSubmissionId, err := primitive.ObjectIDFromHex(input.RatingSubmissionID)
	if err != nil {
		return result, message.RatingSubmissionNotFound
	}

	ratingSubmission, err := s.ratingRepo.GetRatingSubmissionById(ratingSubmissionId)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return result, message.FailedMsg
		}
		return result, message.FailedMsg
	}
	if ratingSubmission == nil {
		return result, message.ErrRatingSubmissionNotFound
	}
	counter = ratingSubmission.LikeCounter

	// Check rating sub helpful exist or not
	ratingSubHelpful, err := s.publicRatingRepo.GetRatingSubHelpfulByRatingSubAndActor(input.RatingSubmissionID, input.UserIDLegacy)
	if err != nil {
		return result, message.FailedMsg
	}
	if ratingSubHelpful == nil {
		counter = counter + 1
		// Create new one of rating sub helpful
		_, errCr := s.publicRatingRepo.CreateRatingSubHelpful(input)
		if errCr != nil {
			if mongo.IsDuplicateKeyError(err) {
				return result, message.ErrDuplicateRatingName
			}
			return result, message.FailedMsg
		}
		// Set result counter & status
		result.Status = "Like"
		result.LikeCounter = counter
	} else {
		if !ratingSubHelpful.Status {
			counter = counter + 1
			result.Status = "Like"
		} else {
			counter = counter - 1
			result.Status = "Unlike"
		}
		// Update status of rating sub helpful
		errUpd := s.publicRatingRepo.UpdateStatusRatingSubHelpful(ratingSubHelpful.ID, ratingSubHelpful.Status)
		if errUpd != nil {
			return result, message.FailedMsg
		}
		// Set result counter & status
		result.LikeCounter = counter
	}

	// Update like_counter rating submission
	err3 := s.publicRatingRepo.UpdateCounterRatingSubmission(ratingSubmissionId, int64(counter))
	if err3 != nil {
		return result, message.ErrSaveData
	}

	return result, message.SuccessMsg
}
