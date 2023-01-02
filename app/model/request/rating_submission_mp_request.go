package request

import (
	"fmt"
	validation "github.com/itgelo/ozzo-validation/v4"
	"go-klikdokter/helper/message"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"strings"
)

// swagger:parameters ListRatingSubmissionMpRequest
type ListRatingSubmissionMpRequest struct {
	Page  int    `json:"page,omitempty" schema:"page" bson:"page"`
	Limit int64  `json:"limit,omitempty" schema:"limit" bson:"limit"`
	Sort  string `json:"sort,omitempty" schema:"sort" bson:"sort"`
	Dir   string `json:"dir,omitempty" schema:"dir" bson:"dir"`
	// Filter available {"user_id_legacy": [""], "source_trans_id": [""], "value": "", "rating_id": [""], "is_with_media": true}
	Filter string `json:"filter"`
}

type RatingSubmissionMpFilter struct {
	UserIDLegacy  []string  `json:"user_uid_legacy"`
	Score         []float64 `json:"score"`
	RatingID      []string  `json:"rating_id"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	SourceTransID string    `json:"source_trans_id"`
	IsWithMedia   *bool     `json:"is_with_media"`
}

// swagger:parameters ReqRatingSubmissionMpBody
type ReqRatingSubmissionMpBody struct {
	// value for product: 20,40,60,80,100
	// value for store: 1,2,3
	// in: body
	Body CreateRatingSubmissionMpRequest `json:"body"`
}

// swagger:parameters ReqRatingSubmissionMpById
type ReqRatingSubmissionMpById struct {
	// ID of Rating Submission
	// in: path
	// required: true
	ID string `json:"id"`
}

// swagger:parameters GetListRatingSummaryMpRequest
type GetListRatingSummaryMpRequest struct {
	// SourceType
	// in: path
	// required: true
	SourceType string `json:"source_type"`
	// Filter available {"source_uid": [], "rating_type": []}
	Filter string `json:"filter" schema:"filter" binding:"omitempty"`
	Limit  int    `json:"limit" schema:"limit" binding:"omitempty,numeric,min=1,max=100"`
	Page   int    `json:"page" schema:"page" binding:"omitempty,numeric,min=1"`
	Sort   string `json:"sort" schema:"sort" binding:"omitempty"`
	Dir    string `json:"dir" schema:"dir" binding:"omitempty"`
}

func (r *GetListRatingSummaryMpRequest) MakeDefaultValueIfEmpty() {
	if r.Limit <= 0 {
		r.Limit = 50
	}
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Sort == "" {
		r.Sort = "updated_at"
	}
}

type CreateRatingSubmissionMpRequest struct {
	UserID        *string        `json:"-" bson:"-"`
	UserIDLegacy  *string        `json:"-" bson:"-"`
	DisplayName   *string        `json:"-" bson:"-"`
	Avatar        string         `json:"-" bson:"-"`
	Value         string         `json:"value" bson:"value"`
	Comment       string         `json:"comment" bson:"comment"`
	IPAddress     string         `json:"ip_address" bson:"ip_address"`
	IsAnonymous   bool           `json:"is_anonymous" bson:"is_anonymous"`
	UserAgent     string         `json:"user_agent" bson:"user_agent"`
	UserPlatform  string         `json:"user_platform" bson:"user_platform"`
	SourceTransID string         `json:"source_trans_id" bson:"source_trans_id"`
	SourceUID     string         `json:"source_uid" bson:"source_uid"`
	RatingType    string         `json:"rating_type" bson:"rating_type"`
	MediaPath     []MediaPathObj `json:"media_path" bson:"media_path"`
}

type MediaPathObj struct {
	UID       string `json:"uid"`
	MediaPath string `json:"media_path"`
}

type RatingSummaryMpResponse struct {
	ID            primitive.ObjectID `json:"id"`
	Name          string             `json:"name,omitempty"`
	Description   *string            `json:"description,omitempty"`
	SourceUid     string             `json:"source_uid,omitempty"`
	SourceType    string             `json:"source_type,omitempty"`
	RatingType    string             `json:"rating_type,omitempty"`
	RatingTypeId  string             `json:"rating_type_id,omitempty"`
	RatingSummary interface{}        `json:"rating_summary,omitempty"`
}

type SaveRatingSubmissionMp struct {
	RatingID      string     `json:"rating_id" bson:"rating_id"`
	UserID        *string    `json:"user_id" bson:"user_id"`
	UserIDLegacy  *string    `json:"user_id_legacy" bson:"user_id_legacy"`
	DisplayName   *string    `json:"display_name" bson:"display_name"`
	Comment       string     `json:"comment" bson:"comment"`
	Value         *string    `json:"value" bson:"value"`
	Avatar        string     `json:"avatar" bson:"avatar"`
	IPAddress     string     `json:"ip_address" bson:"ip_address"`
	UserAgent     string     `json:"user_agent" bson:"user_agent"`
	SourceTransID string     `json:"source_trans_id" bson:"source_trans_id"`
	SourceUID     string     `json:"source_uid" bson:"source_uid"`
	SourceType    string     `json:"source_type" bson:"source_type"`
	UserPlatform  string     `json:"user_platform" bson:"user_platform"`
	Tagging       TaggingObj `json:"tagging" bson:"tagging"`
	IsAnonymous   bool       `json:"is_anonymous" bson:"is_anonymous"`
	MediaPath     []string   `json:"media_path" bson:"media_path"`
	IsWithMedia   bool       `json:"is_with_media" bson:"is_with_media"`
	ReplyComment  string     `json:"reply_comment"`
	ReplyBy       string     `json:"reply_by"`
	OrderNumber   string     `json:"order_number"`
}

func (req CreateRatingSubmissionMpRequest) Validate() error {
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
	)
}

func sliceStringToSliceInterface(arr []string) []interface{} {
	arrInterface := make([]interface{}, len(arr))
	for i, v := range arr {
		arrInterface[i] = v
	}
	return arrInterface
}
