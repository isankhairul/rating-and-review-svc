package request

import (
	validation "github.com/itgelo/ozzo-validation/v4"
	"github.com/spf13/viper"
)

// swagger:parameters CreateRatingRequest
type CreateRatingRequest struct {
	// in: body
	Body SaveRatingRequest `json:"body"`
}

// swagger:parameters UpdateRatingRequest
type UpdateRatingRequest struct {
	// in: path
	// required: true
	Id string `json:"id"`
	// in: body
	Body BodyUpdateRatingRequest `json:"body"`
}

// swagger:model BodyUpdateRatingRequest
type BodyUpdateRatingRequest struct {
	// example: Rumah Sakit - RS Pondok Indah Bintaro Jaya
	Name string `json:"name"`
	// example: Rating Group for Rumah Sakit RS Pondok Indah Bintaro Jaya
	Description *string `json:"description"`
	// example: 2729
	SourceUid string `json:"source_uid"`
	// example: hospital
	SourceType string `json:"source_type"`
	// example: true
	CommentAllowed *bool `json:"comment_allowed"`
}

// swagger:parameters GetListRatingsRequest
type GetListRatingsRequest struct {
	Filter string `json:"filter" schema:"filter" binding:"omitempty"`
	Limit  int    `json:"limit" schema:"limit" binding:"omitempty,numeric,min=1,max=100"`
	Page   int    `json:"page" schema:"page" binding:"omitempty,numeric,min=1"`
	Sort   string `json:"sort" schema:"sort" binding:"omitempty"`
	Dir    string `json:"dir" schema:"dir" binding:"omitempty"`
}

// swagger:model SaveRatingRequest
type SaveRatingRequest struct {
	// example: Rumah Sakit - RS Pondok Indah Bintaro Jaya
	Name string `json:"name"`
	// example: Rating Group for Rumah Sakit RS Pondok Indah Bintaro Jaya
	Description *string `json:"description"`
	// example: 2729
	SourceUid string `json:"source_uid"`
	// example: hospital
	SourceType string `json:"source_type"`
	// example: standard-0.0-to-5.0
	RatingType string `json:"rating_type"`
	// example: 629dc84ff16b9b21f357a609
	RatingTypeId string `json:"rating_type_id"`
	// example: true
	CommentAllowed *bool `json:"comment_allowed"`
	// example: true
	Status *bool `json:"status"`
}

// swagger:parameters GetRatingRequest
type GetRatingRequest struct {
	// in: path
	// required: true
	Id string `json:"id"`
}

// swagger:parameters DeleteRatingRequest
type DeleteRatingRequest struct {
	// in: path
	// required: true
	Id string `json:"id"`
}

func (req SaveRatingRequest) Validate() error {
	sourceType := viper.GetStringSlice("source-type")
	return validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required),
		validation.Field(&req.SourceType, validation.Required, validation.In(sourceType[0], sourceType[1], sourceType[2])),
		validation.Field(&req.SourceUid, validation.Required),
		validation.Field(&req.RatingTypeId, validation.Required),
		validation.Field(&req.RatingType, validation.Required),
	)
}

func (req BodyUpdateRatingRequest) Validate() error {
	sourceType := viper.GetStringSlice("source-type")
	return validation.ValidateStruct(&req,
		validation.Field(&req.SourceType, validation.In(sourceType[0], sourceType[1], sourceType[2])),
	)
}

type RatingFilter struct {
	SourceUid    []string `json:"source_uid"`
	RatingTypeId []string `json:"rating_type_id"`
}

func (r *GetListRatingsRequest) MakeDefaultValueIfEmpty() {
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
