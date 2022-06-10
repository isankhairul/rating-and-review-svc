package request

import (
	validation "github.com/itgelo/ozzo-validation/v4"
	"go-klikdokter/helper/message"
	"regexp"
	"time"
)

var (
	regexIP = "\\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\\.|$)){4}\\b"
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
	UserUID []string `json:"user_uid"`
	Score   []string `json:"score"`
}

// swagger:parameters ReqRatingSubmissonBody
type ReqRatingSubmissonBody struct {
	//  in: body
	Body CreateRatingSubmissonRequest `json:"body"`
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

type CreateRatingSubmissonRequest struct {
	RatingID     string  `json:"rating_id" bson:"rating_id"`
	UserID       *string `json:"user_id" bson:"user_id"`
	UserIDLegacy *string `json:"user_id_legacy" bson:"user_id_legacy"`
	Comment      string  `json:"comment" bson:"comment"`
	Value        float64 `json:"value" bson:"value"`
	IPAddress    string  `json:"ip_address" bson:"ip_address"`
	UserAgent    string  `json:"user_agent" bson:"user_agent"`
}

type UpdateRatingSubmissonRequest struct {
	ID           string    `json:"id"`
	RatingID     string    `json:"rating_id,omitempty" bson:"rating_id"`
	UserID       *string   `json:"user_id,omitempty" bson:"user_id"`
	UserIDLegacy *string   `json:"user_id_legacy,omitempty" bson:"user_id_legacy"`
	Comment      string    `json:"comment,omitempty" bson:"comment"`
	Value        float64   `json:"value,omitempty" bson:"value"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

func (req CreateRatingSubmissonRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.RatingID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Value, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.IPAddress, validation.Match(regexp.MustCompile(regexIP)).Error(message.ErrIPFormatReq.Message)),
	)
}

func (req UpdateRatingSubmissonRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.RatingID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Value, validation.Required.Error(message.ErrReq.Message)),
	)
}
