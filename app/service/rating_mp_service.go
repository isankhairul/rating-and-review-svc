package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-klikdokter/app/middleware"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	publicrequest "go-klikdokter/app/model/request/public"
	"go-klikdokter/app/model/response"
	publicresponse "go-klikdokter/app/model/response/public"
	"go-klikdokter/app/repository"
	helper_dapr "go-klikdokter/helper/dapr"
	global "go-klikdokter/helper/global"
	"go-klikdokter/helper/message"
	"go-klikdokter/helper/thumbor"
	"go-klikdokter/pkg/util"
	util_media "go-klikdokter/pkg/util/media"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/vjeantet/govaluate"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RatingMpService interface {
	// Rating submission
	CreateRatingSubmissionMp(ctx context.Context, input request.CreateRatingSubmissionRequest) ([]response.CreateRatingSubmissionMpResponse, message.Message)
	UpdateRatingSubmission(ctx context.Context, input request.UpdateRatingSubmissionRequest) message.Message
	ReplyAdminRatingSubmission(input request.ReplyAdminRatingSubmissionRequest) message.Message

	// unused
	GetRatingSubmissionMp(id string) (*response.RatingSubmissionMpResponse, message.Message)
	GetListRatingSubmissionsMp(input request.ListRatingSubmissionRequest) ([]response.RatingSubmissionMpResponse, *base.Pagination, message.Message)
	GetListRatingSummaryBySourceType(input request.GetListRatingSummaryRequest) ([]response.RatingSummaryMpResponse, *base.Pagination, message.Message)
	// Rating
	CreateRating(input request.SaveRatingRequest) (*entity.RatingsMpCol, message.Message)
	GetRatingById(id string) (*entity.RatingsMpCol, message.Message)
	UpdateRating(input request.UpdateRatingRequest) message.Message
	DeleteRating(id string) message.Message
	GetListRatings(input request.GetListRatingsRequest) ([]entity.RatingsMpCol, *base.Pagination, message.Message)
	// end unused
}

type ratingMpServiceImpl struct {
	logger       log.Logger
	ratingMpRepo repository.RatingMpRepository
}

func NewRatingMpService(
	lg log.Logger,
	rmp repository.RatingMpRepository,
) RatingMpService {
	return &ratingMpServiceImpl{lg, rmp}
}

func (s *ratingMpServiceImpl) GetRatingById(id string) (*entity.RatingsMpCol, message.Message) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, message.ErrDataNotFound
	}

	result, err := s.ratingMpRepo.GetRatingById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.ErrDataNotFound
		}
		return nil, message.FailedMsg
	}

	return result, message.SuccessMsg
}

func (s *ratingMpServiceImpl) UpdateRating(input request.UpdateRatingRequest) message.Message {
	// get current rating
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return message.ErrDataNotFound
	}

	currentRating, err := s.ratingMpRepo.GetRatingById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return message.ErrDataNotFound
		}
		return message.FailedMsg
	}

	// check rating has rating submission
	ratingSubmission, err := s.ratingMpRepo.GetRatingSubmissionByRatingId(input.Id)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return message.FailedMsg
		}
	}

	if ratingSubmission != nil && (input.Body.SourceUid != "" || input.Body.SourceType != "") {
		return message.ErrCanNotUpdateSourceTypeOrSoureUid
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
		rating, err := s.ratingMpRepo.GetRatingByRatingTypeSourceUidAndSourceType(currentRating.RatingTypeId, sourceUid, sourceType)
		if err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				return message.FailedMsg
			}
		}

		if rating != nil {
			return message.ErrExistingRatingTypeIdSourceUidAndSourceType
		}
	}

	_, err = s.ratingMpRepo.UpdateRating(objectId, input.Body)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return message.ErrDuplicateRatingName
		}
		return message.FailedMsg
	}
	return message.SuccessMsg
}

func (s *ratingMpServiceImpl) CreateRating(input request.SaveRatingRequest) (*entity.RatingsMpCol, message.Message) {

	if input.Status == nil {
		status := true
		input.Status = &status
	}

	// check rating type id, source uid and source type not exist
	rating, err := s.ratingMpRepo.GetRatingByRatingTypeSourceUidAndSourceType(input.RatingTypeId, input.SourceUid, input.SourceType)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.FailedMsg
		}
	}
	if rating != nil {
		return nil, message.ErrExistingRatingTypeIdSourceUidAndSourceType
	}

	// check rating type exist
	ratingTypeId, err := primitive.ObjectIDFromHex(input.RatingTypeId)
	if err != nil {
		return nil, message.ErrRatingTypeNotExist
	}
	ratingTypeNum, err := s.ratingMpRepo.GetRatingTypeNumByIdAndStatus(ratingTypeId)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.FailedMsg
		}
	}

	ratingTypeLikert, err := s.ratingMpRepo.GetRatingTypeLikertByIdAndStatus(ratingTypeId)
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

	result, err := s.ratingMpRepo.CreateRating(input)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, message.ErrDuplicateRatingName
		}
		return nil, message.FailedMsg
	}
	return result, message.SuccessMsg
}

func (s *ratingMpServiceImpl) DeleteRating(id string) message.Message {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return message.ErrDataNotFound
	}
	// FindRatingByRatingID

	rating, err := s.ratingMpRepo.GetRatingById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return message.ErrDataNotFound
		}
		return message.FailedMsg
	}

	ratingSubmission, err := s.ratingMpRepo.GetRatingSubmissionByRatingId(rating.ID.Hex())
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return message.FailedMsg
		}
	}

	if ratingSubmission != nil {
		return message.ErrRatingHasRatingSubmission
	}

	err = s.ratingMpRepo.DeleteRating(objectId)
	if err != nil {
		return message.FailedMsg
	}
	return message.SuccessMsg
}

func (s *ratingMpServiceImpl) GetListRatings(input request.GetListRatingsRequest) ([]entity.RatingsMpCol, *base.Pagination, message.Message) {
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

	ratings, pagination, err := s.ratingMpRepo.GetRatingsByParams(input.Limit, input.Page, dir, input.Sort, filter)
	if err != nil {
		return nil, nil, message.FailedMsg
	}
	results := make([]entity.RatingsMpCol, 0)
	if len(ratings) == 0 {
		return results, pagination, message.SuccessMsg
	}
	results = ratings

	return results, pagination, message.SuccessMsg
}

func (s *ratingMpServiceImpl) CreateRatingSubmissionMp(ctx context.Context, input request.CreateRatingSubmissionRequest) ([]response.CreateRatingSubmissionMpResponse, message.Message) {
	// logger := log.With(s.logger, "RatingService", "CreateRatingSubmission")
	var saveReq = make([]entity.RatingSubmissionMp, 0)
	// set source type
	correlationId := fmt.Sprint(ctx.Value(middleware.CorrelationIdContextKey))
	sourceType := global.GetSourceTypeByRatingType(input.RatingType)
	result := []response.CreateRatingSubmissionMpResponse{}
	// isOrderIdExist := false

	// validation input
	if err := input.ValidateMp(); err != nil {
		return result, message.Message{
			Code:    message.ValidationFailCode,
			Message: err.Error(),
		}
	}
	
	// set user_id as user_id_legacy must be filled
	input.UserID = input.UserIDLegacy
	// Validate displayname
	if input.DisplayName == nil || *input.DisplayName == "" {
		return result, message.ErrDisplayNameRequired
	}

	// The maximum length of user_agent allowed is 200 characters. Crop at 197 characters with triple dots (...) at the end.
	if len(strings.TrimSpace(input.UserAgent)) > 200 {
		return result, message.UserAgentTooLong
	}

	originalSourceTransID := input.SourceTransID

	ratingTypeNum, err := s.ratingMpRepo.FindRatingTypeNumByRatingType(input.RatingType)

	if err != nil {
		return result, message.ErrDB
	}

	if ratingTypeNum == nil {
		return result, message.ErrRatingTypeNotExist
	}

	// Concate source_trans_id, source_type, source_uid, user_id
	input.SourceTransID = originalSourceTransID + "||" + sourceType + "||" + input.SourceUID + "||" + *input.UserID

	// A submission with a combination of either (rating_id and user_id) OR (rating_id and user_id_legacy) is allowed once
	userHasSubmitRating, _ := checkUserHaveSubmitRatingMpBySourceTransID(input.SourceTransID, s)
	if userHasSubmitRating != nil {
		return result, message.UserRated
	}
	// process media_path
	var media []entity.MediaObj
	var isWithMedia bool
	if len(input.Media) > 0 {
		for _, mp := range input.Media {
			if mp.MediaPath != "" && mp.UID != "" {
				media = append(media, entity.MediaObj{UID: mp.UID, MediaPath: mp.MediaPath})
				isWithMedia = true
			}
		}
	}
	// end process media_path
	value, _ := strconv.Atoi(input.Value)
	saveReq = append(saveReq, entity.RatingSubmissionMp{
		// RatingID:      rating.ID.Hex(),
		Value:         value,
		UserID:        input.UserID,
		UserIDLegacy:  input.UserIDLegacy,
		DisplayName:   input.DisplayName,
		Comment:       &input.Comment,
		Avatar:        input.Avatar,
		IPAddress:     input.IPAddress,
		UserAgent:     input.UserAgent,
		SourceTransID: input.SourceTransID,
		UserPlatform:  input.UserPlatform,
		IsAnonymous:   input.IsAnonymous,
		SourceUID:     input.SourceUID,
		SourceType:    sourceType,
		Media:         media,
		IsWithMedia:   isWithMedia,
		OrderNumber:   originalSourceTransID,
		RatingTypeID:  ratingTypeNum.ID.Hex(),
		Cancelled:     false,
		StoreUID:      input.StoreUID,
	})

	if len(saveReq) == 0 {
		return result, message.ErrTypeNotFound
	}

	ratingSubs, err := s.ratingMpRepo.CreateRatingSubmission(saveReq)
	if err != nil {
		return result, message.ErrSaveData
	}

	ratingSubsID := ""

	for _, val := range *ratingSubs {
		ratingSubsID = val.ID.Hex()
	}

	go func() {
		// trigger image house keeping
		util_media.ImageHouseKeeping(ctx, s.logger, media, ratingSubsID)
		// send review for product & store to payment svc
		if ratingSubs != nil && len(*ratingSubs) > 0 {
			ratingSub := *ratingSubs
			util.UpdateReviewProductStore(originalSourceTransID, sourceType, input.SourceUID, ratingSub[0].ID.Hex(), s.logger)
		}
	}()

	if input.RatingType == "rating_for_product" {
		go GetFinalRating(s, correlationId, sourceType, input.SourceUID)
	}
	
	return result, message.SuccessMsg
}

func (s *ratingMpServiceImpl) UpdateRatingSubmission(ctx context.Context ,input request.UpdateRatingSubmissionRequest) message.Message {
	// Input ID of Submission
	objectRatingSubmissionId, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return message.RatingSubmissionNotFound
	}
	// find Rating submission
	ratingSubmission, err := s.ratingMpRepo.GetRatingSubmissionById(objectRatingSubmissionId)
	if err != nil || ratingSubmission == nil {
		return message.RatingSubmissionNotFound
	}

	// Validate value of numeric type
	var validateMsg message.Message

	if validateMsg.Code == message.ValidationFailCode {
		return validateMsg
	}

	// validate cannot update rating submission of another user
	notValidUpdate := util.ValidateUserCannotUpdateMp(*input.UserID, *input.UserIDLegacy, *ratingSubmission)
	if notValidUpdate {
		return message.ErrUserPermissionUpdate
	}

	// set update data ratingSub
	var timeUpdate time.Time
	timeUpdate = time.Now().In(util.Loc)
	// process media_path
	var media []entity.MediaObj
	var isWithMedia bool
	if len(input.Media) > 0 {
		for _, mp := range input.Media {
			if mp.MediaPath != "" && mp.UID != "" {
				media = append(media, entity.MediaObj{UID: mp.UID, MediaPath: mp.MediaPath})
				isWithMedia = true
			}
		}
	}
	// end process media_path
	value, _ := strconv.Atoi(*input.Value)
	ratingSubmission.Comment = &input.Comment
	ratingSubmission.Value = value
	ratingSubmission.Media = media
	ratingSubmission.IsWithMedia = isWithMedia
	ratingSubmission.UpdatedAt = timeUpdate

	// Update
	errC := s.ratingMpRepo.UpdateRatingSubmission(*ratingSubmission, objectRatingSubmissionId)
	if errC != nil {
		return message.ErrSaveData
	}

	go func() {
		// trigger image house keeping
		util_media.ImageHouseKeeping(ctx, s.logger, media, ratingSubmission.ID.Hex())
	}()

	return message.SuccessMsg

}

// swagger:route PUT /rating-submissions/reply/{id} RatingSubmission ReqAdminReplyRatingSubmissionBody
// Reply Rating Submission By ID
//
// security:
// - Bearer: []
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ratingMpServiceImpl) ReplyAdminRatingSubmission(input request.ReplyAdminRatingSubmissionRequest) message.Message {
	// Input ID of Submission
	objectRatingSubmissionId, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return message.RatingSubmissionNotFound
	}
	// find Rating submission
	ratingSubmission, err := s.ratingMpRepo.GetRatingSubmissionById(objectRatingSubmissionId)
	if err != nil || ratingSubmission == nil {
		return message.RatingSubmissionNotFound
	}

	// set update data ratingSub
	ratingSubmission.UpdatedAt = time.Now().In(util.Loc)
	ratingSubmission.Reply = input.Reply
	ratingSubmission.ReplyBy = fmt.Sprint(input.JWTObj.Fullname)

	// Update
	errC := s.ratingMpRepo.UpdateRatingSubmission(*ratingSubmission, objectRatingSubmissionId)
	if errC != nil {
		return message.ErrSaveData
	}

	return message.SuccessMsg
}

func (s *ratingMpServiceImpl) GetRatingSubmissionMp(id string) (*response.RatingSubmissionMpResponse, message.Message) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, message.ErrRatingSubmissionNotFound
	}
	get, err := s.ratingMpRepo.GetRatingSubmissionById(objectId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, message.ErrRatingSubmissionNotFound
		}
		return nil, message.FailedMsg
	}
	var mediaResponse []response.MediaObjResponse
	for _, value := range get.Media {
		mediaObjResponse := response.MediaObjResponse{
			UID:        value.UID,
			MediaPath:  value.MediaPath,
			MediaImage: thumbor.GetNewThumborImagesOriginal(value.MediaPath),
		}
		mediaResponse = append(mediaResponse, mediaObjResponse)
	}
	
	var result = response.RatingSubmissionMpResponse{
		RatingID:      get.RatingID,
		UserID:        get.UserID,
		UserIDLegacy:  get.UserIDLegacy,
		Value:         strconv.Itoa(get.Value),
		SourceTransID: get.SourceTransID,
		Media:         mediaResponse,
		IsWithMedia:   get.IsWithMedia,
	}
	if get != nil && get.Comment != nil {
		result.Comment = *get.Comment
	}

	return &result, message.SuccessMsg
}

func (s *ratingMpServiceImpl) GetListRatingSubmissionsMp(input request.ListRatingSubmissionRequest) ([]response.RatingSubmissionMpResponse, *base.Pagination, message.Message) {
	var dir interface{}
	results := make([]response.RatingSubmissionMpResponse, 0)
	userIdEmpty := ""
	commentEmpty := ""
	if input.Dir == "asc" {
		dir = 1
	} else {
		dir = -1
	}
	// Set default value
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 50
	}
	if input.Sort == "" {
		input.Sort = "created_at"
	}

	filter := request.RatingSubmissionMpFilter{}
	if input.Filter != "" {
		errMarshal := json.Unmarshal([]byte(input.Filter), &filter)
		if errMarshal != nil {
			return nil, nil, message.WrongFilter
		}
	}
	ratingSubmissions, pagination, err := s.ratingMpRepo.GetListRatingSubmissions(filter, input.Page, input.Limit, input.Sort, dir)
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
		if filScore := filterScoreSubmissionMp(args, filter.Score); filScore {
			// create thumbor response
			var mediaResponse []response.MediaObjResponse
			for _, value := range args.Media {
				mediaObjResponse := response.MediaObjResponse{
					UID:        value.UID,
					MediaPath:  value.MediaPath,
					MediaImage: thumbor.GetNewThumborImagesOriginal(value.MediaPath),
				}
				mediaResponse = append(mediaResponse, mediaObjResponse)
			}

			results = append(results, response.RatingSubmissionMpResponse{
				RatingID:      args.RatingID,
				UserID:        args.UserID,
				UserIDLegacy:  args.UserIDLegacy,
				Comment:       *args.Comment,
				Value:         strconv.Itoa(args.Value),
				SourceTransID: args.SourceTransID,
				Media:         mediaResponse,
				IsWithMedia:   args.IsWithMedia,
			})
		}
	}

	if len(filter.Score) > 0 && pagination != nil {
		pagination.Records = int64(len(results))
		pagination.TotalRecords = int64(len(results))
	}

	return results, pagination, message.SuccessMsg
}

func checkUserHaveSubmitRatingMp(userId, ratingId, sourceTransId string, s *ratingMpServiceImpl) (*entity.RatingSubmissionMp, error) {

	ratingSubmissionMp, err := s.ratingMpRepo.FindRatingSubmissionByUserIDAndRatingID(&userId, ratingId, sourceTransId)
	return ratingSubmissionMp, err
}

func checkUserHaveSubmitRatingMpBySourceTransID(sourceTransId string, s *ratingMpServiceImpl) (*entity.RatingSubmissionMp, error) {

	ratingSubmissionMp, err := s.ratingMpRepo.FindRatingSubmissionBySourceTransID(sourceTransId)
	return ratingSubmissionMp, err
}

func filterScoreSubmissionMp(ratingSubmissionsMp entity.RatingSubmissionMp, score []float64) bool {
	if len(score) == 0 {
		return true
	}
	var scoreDB = make([]float64, 0)
	
	values := strings.Split(strconv.Itoa(ratingSubmissionsMp.Value), ",")
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

func (s *ratingMpServiceImpl) GetListRatingSummaryBySourceType(input request.GetListRatingSummaryRequest) ([]response.RatingSummaryMpResponse, *base.Pagination, message.Message) {
	results := []response.RatingSummaryMpResponse{}
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

	ratings, pagination, err := s.ratingMpRepo.GetPublicRatingsByParams(input.Limit, input.Page, dir, input.Sort, filter)
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

func (s *ratingMpServiceImpl) summaryRatingLikert(rating entity.RatingsMpCol, ratingLikert entity.RatingTypesLikertCol) (*response.RatingSummaryMpResponse, error) {
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
			totalCount, err := s.ratingMpRepo.CountRatingSubsByRatingIdAndValue(rating.ID.Hex(), strconv.Itoa(i))
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

	result := response.RatingSummaryMpResponse{
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

func (s *ratingMpServiceImpl) summaryRatingNumeric(rating entity.RatingsMpCol, sourceType string) (*response.RatingSummaryMpResponse, error) {
	sumCountRatingSubs, err := s.ratingMpRepo.GetSumCountRatingSubsByRatingId(rating.ID.Hex())
	if err != nil {
		return nil, err
	}
	if sumCountRatingSubs == nil {
		return nil, errors.New("data RatingSubmission not found")
	}
	formulaRating, err := s.ratingMpRepo.GetRatingFormulaByRatingTypeIdAndSourceType(rating.RatingTypeId, sourceType)
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

		data := response.RatingSummaryMpResponse{
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
	result.TotalValue = "0"

	for _, c := range sumCountRatingSubs.Comments {
		if strings.TrimSpace(c) != "" {
			result.TotalComment++
		}
	}

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
		// totalValue, err := strconv.ParseFloat(fmt.Sprintf("%.1f", finalCalc.(float64)), 64)
		// if err != nil {
		//	return result, err
		// }
		result.TotalValue = fmt.Sprintf("%.1f", finalCalc)
	}

	return result, nil
}

func GetFinalRating(s *ratingMpServiceImpl, correlationId string ,sourceType string, sourceUid string) {
	if s.logger == nil {
		return
	}
	summarySubsGroupByValue, _ := s.ratingMpRepo.GetRatingSubsGroupByValue(sourceUid, sourceType)
	totalValue := 0
	TotalReviewer := 0

	for _, val := range summarySubsGroupByValue {
		totalValue += val.ConvertedValue * val.Total
		TotalReviewer += val.Total
	}

	formulaRating, _  := s.ratingMpRepo.GetRatingFormulaBySourceType(sourceType)
	
	if formulaRating == nil {
		level.Info(s.logger).Log("correlationID", correlationId, "Type", "Create Rating Submission", "err", "Failed Get Formula Rating")	
		return
	}

	formula := formulaRating.Formula
	expression, err := govaluate.NewEvaluableExpression(formula)

	if err != nil {
		level.Info(s.logger).Log("correlationID", correlationId, "Type", "Create Rating Submission", "err", "Failed Get Expression Formula Rating")
		return
	}

	parameters := make(map[string]interface{}, 8)
	parameters["sum"] = totalValue
	parameters["count"] = TotalReviewer

	finalCalc, err := expression.Evaluate(parameters)
	if err != nil {
		level.Info(s.logger).Log("correlationID", correlationId, "Type", "Create Rating Submission", "err", "Failed Calculate Formula Rating")
		return
	}
	result := fmt.Sprintf("%.1f", finalCalc)
	level.Info(s.logger).Log("correlationID", correlationId, "Type", "Create Rating Submission","result", result)

	jsonPayload, _ := json.Marshal(map[string]interface{}{
		"data": map[string]string{
			"product_uid": sourceUid,
			"final_rating": result,
		},
	})
		
	client := helper_dapr.NewDaprHttpClient()
	response, err := client.PublishEvent("queuing.rnr.ts-final-rating", string(jsonPayload))

	_ = level.Info(s.logger).Log(fmt.Sprintf("Publisher response: %v", response))
	_ = level.Info(s.logger).Log(fmt.Sprintf("Publisher error: %v", err))
}