package request

import (
	validation "github.com/itgelo/ozzo-validation/v4"
	"github.com/spf13/viper"
)

// swagger:parameters GetPublicListRatingSubmissionByUserIdRequest
type GetPublicListRatingSubmissionByUserIdRequest struct {
	// SourceType
	// in: path
	// required: true
	SourceType string `json:"source_type"`

	// SourceUID
	// in: path
	// required: true
	SourceUID string `json:"source_uid"`

	// UserIdLegacy
	// in: path
	// required: true
	UserIdLegacy string `json:"user_id_legacy"`

	Limit int    `json:"limit" schema:"limit" binding:"omitempty,numeric,min=1,max=100"`
	Page  int    `json:"page" schema:"page" binding:"omitempty,numeric,min=1"`
	Sort  string `json:"sort" schema:"sort" binding:"omitempty"`
	Dir   string `json:"dir" schema:"dir" binding:"omitempty"`
}

// swagger:parameters ReqUpdateRatingSubDisplayNameBody
type ReqUpdateRatingSubDisplayNameBody struct {
	// in: path
	// required: true
	UserIdLegacy string `json:"user_id_legacy"`
	// in: body
	Body UpdateRatingSubDisplayNameRequest `json:"body"`
}

type UpdateRatingSubDisplayNameRequest struct {
	UserIdLegacy string `json:"-"`
	DisplayName  string `json:"display_name,omitempty"`
}

func (req GetPublicListRatingSubmissionByUserIdRequest) ValidateSourceType() error {
	sourceType := viper.GetStringSlice("source-type")
	return validation.ValidateStruct(&req,
		validation.Field(&req.SourceType, validation.In(sourceType[0], sourceType[1], sourceType[2])),
	)
}
