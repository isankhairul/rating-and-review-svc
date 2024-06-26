package publicrequest

import (
	"fmt"
	"strings"

	validation "github.com/itgelo/ozzo-validation/v4"
	"github.com/spf13/viper"
)

// swagger:parameters GetPublicListRatingSummaryMpRequest
type GetPublicListRatingSummaryMpRequest struct {
	// SourceType
	// in: path
	// required: true
	SourceType string `json:"source_type"`
	// Filter available {"source_uid": [""], "rating_type": [""]}
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

// swagger:parameters GetPublicListRatingSubmissionByIDRequest
type GetPublicListRatingSubmissionByIDRequest struct {
	// Filter available {"rating_subs_id": [""]}
	Filter string `json:"filter" schema:"filter" binding:"omitempty"`
	Source string `json:"source" schema:"source" binding:"omitempty"`
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

func (req GetPublicListRatingSubmissionByIDRequest) ValidateFilterAndSource() error {
	listAvailableSource := []string{"all", "mp"}

	return validation.ValidateStruct(&req,
		validation.Field(&req.Source, validation.Required, validation.In(sliceStringToSliceInterface(listAvailableSource)...).Error(fmt.Sprintf("source should be %s", strings.Join(listAvailableSource, ",")))),
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
	RatingSubsID  []string     `json:"rating_subs_id"`
	Value         string       `json:"value"`
	IsWithMedia   *bool        `json:"is_with_media"`
	StartDate     string       `json:"start_date"`
	EndDate       string       `json:"end_date"`
}

func (req FilterRatingSubmissionMp) ValidateFormatDate() error {
	return validation.ValidateStruct(&req,
		// Default
		validation.Field(&req.StartDate, validation.When(req.StartDate != "", validation.Date("2006-01-02").
			Error("start_date is invalid format, format should be 2006-01-02"))),
		validation.Field(&req.EndDate, validation.When(req.EndDate != "", validation.Date("2006-01-02").
			Error("end_date is invalid format, format should be 2006-01-02"))),
	)
}

func sliceStringToSliceInterface(arr []string) []interface{} {
	arrInterface := make([]interface{}, len(arr))
	for i, v := range arr {
		arrInterface[i] = v
	}
	return arrInterface
}

// swagger:parameters PublicGetListDetailRatingSummaryRequest
type PublicGetListDetailRatingSummaryRequest struct {
	// Filter available {"source_uid": [""]}
	// max source_uid is 50
	// required: true
	Filter string `json:"filter" schema:"filter" binding:"omitempty"`

	// Limit int    `json:"limit" schema:"limit" binding:"omitempty,numeric,min=1,max=100"`
	// Page  int    `json:"page" schema:"page" binding:"omitempty,numeric,min=1"`
	// Sort  string `json:"sort" schema:"sort" binding:"omitempty"`
	// Dir   string `json:"dir" schema:"dir" binding:"omitempty"`

	// SourceType
	// in: path
	// required: true
	SourceType string `json:"source_type"`
}

func (r *PublicGetListDetailRatingSummaryRequest) MakeDefaultValueIfEmpty() {
	// if r.Limit <= 0 {
	// 	r.Limit = 50
	// }
	// if r.Page <= 0 {
	// 	r.Page = 1
	// }
	// if r.Sort == "" {
	// 	r.Sort = "updated_at"
	// }
}

// swagger:parameters PublicGetRatingSummaryStoreProductRequest
type PublicGetRatingSummaryStoreProductRequest struct {
	// Filter available {"store_uid": [""]}
	// max store_uid is 20
	// required: true
	Filter string `json:"filter" schema:"filter" binding:"omitempty"`

	// Limit int    `json:"-"`
	// Page  int    `json:"-"`
	// Sort string `json:"sort" schema:"sort" binding:"omitempty"`
	// Dir  string `json:"dir" schema:"dir" binding:"omitempty"`
}

func (r *PublicGetRatingSummaryStoreProductRequest) MakeDefaultValueIfEmpty() {
	// if r.Limit <= 0 {
	// 	r.Limit = 10
	// }
	// if r.Page <= 0 {
	// 	r.Page = 1
	// }
	// if r.Sort == "" {
	// 	r.Sort = "updated_at"
	// }
}
