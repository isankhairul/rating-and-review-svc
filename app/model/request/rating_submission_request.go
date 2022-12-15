package request

import (
	"go-klikdokter/helper/message"
	"regexp"
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
	Ratings       []RatingByType `json:"ratings" bson:"ratings"`
	UserID        *string        `json:"user_id" bson:"user_id"`
	UserIDLegacy  *string        `json:"user_id_legacy" bson:"user_id_legacy"`
	DisplayName   *string        `json:"display_name" bson:"display_name"`
	Comment       string         `json:"comment" bson:"comment"`
	IPAddress     string         `json:"ip_address" bson:"ip_address"`
	UserAgent     string         `json:"user_agent" bson:"user_agent"`
	SourceTransID string         `json:"source_trans_id" bson:"source_trans_id"`
	UserPlatform  string         `json:"user_platform" bson:"user_platform"`
	Avatar        string         `json:"-" bson:"-"`
	IsAnonymous   bool           `json:"is_anonymous" bson:"is_anonymous"`
	SourceUID     string         `json:"source_uid" bson:"source_uid"`
	RatingType    string         `json:"rating_type" bson:"rating_type"`
	Value         string         `json:"value" bson:"value"`
	MediaPath     []MediaPathObj `json:"media_path" bson:"media_path"`
	Token         string         `json:"-"`
}

type SaveRatingSubmission struct {
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
}

type TaggingObj struct {
	RatingId string   `json:"rating_id" bson:"rating_id"`
	Value    []string `json:"value" bson:"value"`
}

type UpdateRatingSubmissionRequest struct {
	ID        string    `json:"-"`
	RatingID  string    `json:"rating_id,omitempty" bson:"rating_id"`
	Comment   string    `json:"comment,omitempty" bson:"comment"`
	Value     *string   `json:"value,omitempty" bson:"value"`
	UpdatedAt time.Time `json:"-,omitempty" bson:"updated_at"`
}

func (req CreateRatingSubmissionRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.IPAddress, validation.Match(regexp.MustCompile(regexIP)).Error(message.ErrIPFormatReq.Message)),
		validation.Field(&req.SourceTransID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Comment, validation.NotNil),
	)
}

func (req RatingByType) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Value, validation.Required.Error(message.ErrReq.Message)),
	)
}

func (req UpdateRatingSubmissionRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.RatingID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Value, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Comment, validation.NotNil),
	)
}
