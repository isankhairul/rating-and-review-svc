package publicrequest

import (
	"fmt"
	validation "github.com/itgelo/ozzo-validation/v4"
	"github.com/spf13/viper"
	"strings"
)

// swagger:parameters GetPublicListRatingSummaryMpRequest
type GetPublicListRatingSummaryMpRequest struct {
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

// swagger:parameters GetPublicListRatingSubmissionMpRequest
type GetPublicListRatingSubmissionMpRequest struct {
	// SourceType
	// in: path
	// required: true
	SourceType string `json:"source_type"`
	// SourceUID
	// in: path
	// required: true
	SourceUID string `json:"source_uid"`
	// Filter available {"user_id_legacy": [""], "source_trans_id": [""], "value": "", "is_with_media": true}
	Filter string `json:"filter" schema:"filter" binding:"omitempty"`
	Limit  int    `json:"limit" schema:"limit" binding:"omitempty,numeric,min=1,max=100"`
	Page   int    `json:"page" schema:"page" binding:"omitempty,numeric,min=1"`
	Sort   string `json:"sort" schema:"sort" binding:"omitempty"`
	Dir    string `json:"dir" schema:"dir" binding:"omitempty"`
}

func (r *GetPublicListRatingSummaryMpRequest) MakeDefaultValueIfEmpty() {
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

func (req GetPublicListRatingSubmissionMpRequest) ValidateSourceType() error {
	sourceType := viper.GetStringSlice("source-type")
	interfaceAllSource := make([]interface{}, len(sourceType))
	for i, v := range sourceType {
		interfaceAllSource[i] = v
	}
	return validation.ValidateStruct(&req,
		validation.Field(&req.SourceType, validation.Required, validation.In(interfaceAllSource...).Error(fmt.Sprintf("source_type should be %s", strings.Join(sourceType, ",")))),
	)
}

type FilterRatingSubmissionMp struct {
	SourceUID     string       `json:"source_uid"`
	SourceUIDs    []string     `json:"source_uids"`
	SourceType    string       `json:"source_type"`
	RatingID      []string     `json:"rating_id"`
	LikertFilter  LikertFilter `json:"likert_filter"`
	UserIdLegacy  []string     `json:"user_id_legacy"`
	SourceTransID []string     `json:"source_trans_id"`
	Value         string       `json:"value"`
	IsWithMedia   *bool        `json:"is_with_media"`
}
