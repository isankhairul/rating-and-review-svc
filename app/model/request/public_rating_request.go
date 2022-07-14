package request

import (
	validation "github.com/itgelo/ozzo-validation/v4"
	"github.com/spf13/viper"
)

// swagger:parameters GetPublicListRatingSummaryRequest
type GetPublicListRatingSummaryRequest struct {
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

type FilterRatingSummary struct {
	SourceType string   `json:"source_type"`
	SourceUid  []string `json:"source_uid"`
	RatingType []string `json:"rating_type"`
}

// swagger:parameters GetPublicListRatingSubmissionRequest
type GetPublicListRatingSubmissionRequest struct {
	// SourceType
	// in: path
	// required: true
	SourceType string `json:"source_type"`
	// SourceUID
	// in: path
	// required: true
	SourceUID string `json:"source_uid"`
	Limit     int    `json:"limit" schema:"limit" binding:"omitempty,numeric,min=1,max=100"`
	Page      int    `json:"page" schema:"page" binding:"omitempty,numeric,min=1"`
	Sort      string `json:"sort" schema:"sort" binding:"omitempty"`
	Dir       string `json:"dir" schema:"dir" binding:"omitempty"`
}

type FilterRatingSubmission struct {
	RatingID []string `json:"rating_id"`
}

// swagger:parameters GetRatingBySourceTypeAndActor
type GetRatingBySourceTypeAndActorRequest struct {
	// in: path
	// required: true
	SourceType string `json:"source_type"`
	// in: path
	// required: true
	SourceUID string `json:"source_uid"`
}

func (r *GetPublicListRatingSummaryRequest) MakeDefaultValueIfEmpty() {
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

func (req GetPublicListRatingSubmissionRequest) ValidateSourceType() error {
	sourceType := viper.GetStringSlice("source-type")
	return validation.ValidateStruct(&req,
		validation.Field(&req.SourceType, validation.In(sourceType[0], sourceType[1], sourceType[2])),
	)
}
