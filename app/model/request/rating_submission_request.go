package request

import (
	validation "github.com/itgelo/ozzo-validation/v4"
	"go-klikdokter/helper/message"
	"regexp"
	"time"
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
	UserID    []string `json:"user_uid"`
	Score     []string `json:"score"`
	RatingID  []string `json:"rating_id"`
	StartDate string   `json:"start_date"`
	EndDate   string   `json:"end_date"`
}

// swagger:parameters ReqRatingSubmissonBody
type ReqRatingSubmissonBody struct {
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

// swagger:parameters ReqUpdateRatingSubmissonBody
type ReqUpdateRatingSubmissonBody struct {
	// ID of Rating Submission
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body UpdateRatingSubmissonRequest `json:"body"`
}

type RatingByType struct {
	ID    string  `json:"uid" bson:"uid"`
	Value *string `json:"value" bson:"value"`
}

type CreateRatingSubmissionRequest struct {
	Ratings      []RatingByType `json:"ratings" bson:"ratings"`
	UserID       *string        `json:"user_id" bson:"user_id"`
	UserIDLegacy *string        `json:"user_id_legacy" bson:"user_id_legacy"`
	Comment      string         `json:"comment" bson:"comment"`
	IPAddress    string         `json:"ip_address" bson:"ip_address"`
	UserAgent    string         `json:"user_agent" bson:"user_agent"`
}

type SaveRatingSubmission struct {
	RatingID     string  `json:"rating_id" bson:"rating_id"`
	UserID       *string `json:"user_id" bson:"user_id"`
	UserIDLegacy *string `json:"user_id_legacy" bson:"user_id_legacy"`
	Comment      string  `json:"comment" bson:"comment"`
	Value        *string `json:"value" bson:"value"`
	IPAddress    string  `json:"ip_address" bson:"ip_address"`
	UserAgent    string  `json:"user_agent" bson:"user_agent"`
}

type UpdateRatingSubmissonRequest struct {
	ID           string    `json:"-"`
	RatingID     string    `json:"rating_id,omitempty" bson:"rating_id"`
	UserID       *string   `json:"user_id,omitempty" bson:"user_id"`
	UserIDLegacy *string   `json:"user_id_legacy,omitempty" bson:"user_id_legacy"`
	Comment      string    `json:"comment,omitempty" bson:"comment"`
	Value        *string   `json:"value,omitempty" bson:"value"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

func (req CreateRatingSubmissionRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Ratings, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.IPAddress, validation.Match(regexp.MustCompile(regexIP)).Error(message.ErrIPFormatReq.Message)),
	)
}

func (req RatingByType) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Value, validation.Required.Error(message.ErrReq.Message)),
	)
}

func (req UpdateRatingSubmissonRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.RatingID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Value, validation.Required.Error(message.ErrReq.Message)),
	)
}
