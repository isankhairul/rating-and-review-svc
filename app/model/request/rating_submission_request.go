package request

import (
	"fmt"
	"github.com/spf13/viper"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"
	"regexp"
	"strings"
	"time"

	validation "github.com/itgelo/ozzo-validation/v4"
)

var (
	regexIP = "^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$"
)

// swagger:parameters ListRatingSubmissionRequest
type ListRatingSubmissionRequest struct {
	Page   int    `json:"page,omitempty" schema:"page" bson:"page"`
	Limit  int64  `json:"limit,omitempty" schema:"limit" bson:"limit"`
	Sort   string `json:"sort,omitempty" schema:"sort" bson:"sort"`
	Dir    string `json:"dir,omitempty" schema:"dir" bson:"dir"`
	Filter string `json:"filter"`
}

type RatingSubmissionFilter struct {
	UserIDLegacy  []string  `json:"user_uid_legacy"`
	Score         []float64 `json:"score"`
	RatingID      []string  `json:"rating_id"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	SourceTransID string    `json:"source_trans_id"`
}

// swagger:parameters ReqRatingSubmissionBody ReqPublicRatingSubmissionBody
type ReqRatingSubmissionBody struct {
	//  in: body
	Body CreateRatingSubmissionRequest `json:"body"`
}

// swagger:parameters ReqRatingSubmissionById ReqDeleteRatingSubmissionById
type ReqRatingSubmissionById struct {
	// ID of Rating Submission
	// in: path
	// required: true
	ID string `json:"id"`
}

// swagger:parameters ReqUpdateRatingSubmissionBody
type ReqUpdateRatingSubmissionBody struct {
	// ID of Rating Submission
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body UpdateRatingSubmissionRequest `json:"body"`
}

// swagger:parameters ReqAdminReplyRatingSubmissionBody
type ReqAdminReplyRatingSubmissionBody struct {
	// ID of Rating Submission
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body ReplyAdminRatingSubmissionRequest `json:"body"`
}

// swagger:parameters ReqCancelRatingSubmission
type ReqCancelRatingSubmission struct {
	// in: body
	// required: true
	Body CancelRatingById `json:"body"`
}

type CancelRatingById struct {
	RatingSubmissionId []string `json:"rating_submission_id" bson:"rating_submission_id"`
	CancelledReason    string   `json:"cancelled_reason" bson:"cancelled_reason"`
}

type RatingByType struct {
	ID    string  `json:"uid" bson:"uid"`
	Value *string `json:"value" bson:"value"`
}

type CreateRatingSubmissionRequest struct {
	Ratings       []RatingByType    `json:"ratings" bson:"ratings"`
	UserID        *string           `json:"user_id" bson:"user_id"`
	UserIDLegacy  *string           `json:"user_id_legacy" bson:"user_id_legacy"`
	DisplayName   *string           `json:"display_name" bson:"display_name"`
	Comment       string            `json:"comment" bson:"comment"`
	IPAddress     string            `json:"ip_address" bson:"ip_address"`
	UserAgent     string            `json:"user_agent" bson:"user_agent"`
	SourceTransID string            `json:"source_trans_id" bson:"source_trans_id"`
	UserPlatform  string            `json:"user_platform" bson:"user_platform"`
	Avatar        string            `json:"-" bson:"-"`
	IsAnonymous   bool              `json:"is_anonymous" bson:"is_anonymous"`
	SourceUID     string            `json:"source_uid" bson:"source_uid"`
	StoreUID      string            `json:"store_uid" bson:"store_uid"`
	RatingType    string            `json:"rating_type" bson:"rating_type"`
	Value         string            `json:"value" bson:"value"`
	Media         []entity.MediaObj `json:"media" bson:"media"`
}

type SaveRatingSubmission struct {
	RatingID      string            `json:"rating_id" bson:"rating_id"`
	UserID        *string           `json:"user_id" bson:"user_id"`
	UserIDLegacy  *string           `json:"user_id_legacy" bson:"user_id_legacy"`
	DisplayName   *string           `json:"display_name" bson:"display_name"`
	Comment       string            `json:"comment" bson:"comment"`
	Value         *string           `json:"value" bson:"value"`
	Avatar        string            `json:"avatar" bson:"avatar"`
	IPAddress     string            `json:"ip_address" bson:"ip_address"`
	UserAgent     string            `json:"user_agent" bson:"user_agent"`
	SourceTransID string            `json:"source_trans_id" bson:"source_trans_id"`
	SourceUID     string            `json:"source_uid" bson:"source_uid"`
	SourceType    string            `json:"source_type" bson:"source_type"`
	UserPlatform  string            `json:"user_platform" bson:"user_platform"`
	Tagging       TaggingObj        `json:"tagging" bson:"tagging"`
	IsAnonymous   bool              `json:"is_anonymous" bson:"is_anonymous"`
	Media         []entity.MediaObj `json:"media" bson:"media"`
	IsWithMedia   bool              `json:"is_with_media" bson:"is_with_media"`
}

type TaggingObj struct {
	RatingId string   `json:"rating_id" bson:"rating_id"`
	Value    []string `json:"value" bson:"value"`
}

type ReplyAdminRatingSubmissionRequest struct {
	ID     string        `json:"-"`
	Source string        `json:"source"`
	Reply  string        `json:"reply"`
	JWTObj global.JWTObj `json:"-"`
}

func (req ReplyAdminRatingSubmissionRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Reply, validation.Required.Error(message.ErrReq.Message)),
	)
}

type UpdateRatingSubmissionRequest struct {
	ID           string            `json:"-"`
	RatingType   string            `json:"rating_type"`
	RatingID     string            `json:"rating_id,omitempty" bson:"rating_id"`
	Comment      string            `json:"comment,omitempty" bson:"comment"`
	Value        *string           `json:"value,omitempty" bson:"value"`
	Media        []entity.MediaObj `json:"media" bson:"media"`
	UpdatedAt    time.Time         `json:"-,omitempty" bson:"updated_at"`
	UserID       *string           `json:"-" bson:"user_id"`
	UserIDLegacy *string           `json:"-" bson:"user_id_legacy"`
}

func (req CreateRatingSubmissionRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.IPAddress, validation.Match(regexp.MustCompile(regexIP)).Error(message.ErrIPFormatReq.Message)),
		validation.Field(&req.SourceTransID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Comment, validation.NotNil),
		validation.Field(&req.UserIDLegacy, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.DisplayName, validation.Required.Error(message.ErrReq.Message)),
	)
}

func (req CreateRatingSubmissionRequest) ValidateMp() error {
	arrAllowedValueRatingProduct := []string{"1", "2", "3", "4", "5"}
	arrAllowedValueRatingStore := []string{"1", "2", "3"}

	return validation.ValidateStruct(&req,
		validation.Field(&req.Value, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.SourceUID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.SourceTransID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.SourceTransID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.IPAddress, validation.Match(regexp.MustCompile(regexIP)).Error(message.ErrIPFormatReq.Message)),
		validation.Field(&req.Comment, validation.NotNil),
		validation.Field(&req.Value, validation.When(req.RatingType == "rating_for_product",
			validation.In(sliceStringToSliceInterface(arrAllowedValueRatingProduct)...).Error(fmt.Sprintf("value should be %s", strings.Join(arrAllowedValueRatingProduct, ","))))),
		validation.Field(&req.Value, validation.When(req.RatingType == "rating_for_store",
			validation.In(sliceStringToSliceInterface(arrAllowedValueRatingStore)...).Error(fmt.Sprintf("value should be %s", strings.Join(arrAllowedValueRatingStore, ","))))),
		validation.Field(&req.StoreUID, validation.When(req.RatingType == "rating_for_product", validation.Required)),
	)
}

func (req RatingByType) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Value, validation.Required.Error(message.ErrReq.Message)),
	)
}

func (req UpdateRatingSubmissionRequest) Validate() error {
	arrRatingTypeMp := viper.GetStringSlice("rating-type-mp")
	isRatingMp := stringInSlice(req.RatingType, arrRatingTypeMp)
	arrAllowedValueRatingProduct := []string{"1", "2", "3", "4", "5"}
	arrAllowedValueRatingStore := []string{"1", "2", "3"}

	return validation.ValidateStruct(&req,
		// validation.Field(&req.RatingID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.RatingID, validation.When(isRatingMp == false, validation.Required.Error(message.ErrReq.Message))),
		validation.Field(&req.Value, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Comment, validation.NotNil),
		validation.Field(&req.Value, validation.When(req.RatingType == "rating_for_product",
			validation.In(sliceStringToSliceInterface(arrAllowedValueRatingProduct)...).Error(fmt.Sprintf("value should be %s", strings.Join(arrAllowedValueRatingProduct, ","))))),
		validation.Field(&req.Value, validation.When(req.RatingType == "rating_for_store",
			validation.In(sliceStringToSliceInterface(arrAllowedValueRatingStore)...).Error(fmt.Sprintf("value should be %s", strings.Join(arrAllowedValueRatingStore, ","))))),
	)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
