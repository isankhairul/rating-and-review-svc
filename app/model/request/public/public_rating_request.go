package publicrequest

import (
	"fmt"
	validation "github.com/itgelo/ozzo-validation/v4"
	"github.com/spf13/viper"
	"go-klikdokter/helper/message"
	"strings"
)

// swagger:parameters GetPublicListRatingSummaryRequest
type GetPublicListRatingSummaryRequest struct {
	// SourceType
	// in: path
	// required: true
	SourceType string `json:"source_type"`
	// Filter available {"source_uid": [""]}
	Filter string `json:"filter" schema:"filter" binding:"omitempty"`
	Limit  int    `json:"limit" schema:"limit" binding:"omitempty,numeric,min=1,max=100"`
	Page   int    `json:"page" schema:"page" binding:"omitempty,numeric,min=1"`
	Sort   string `json:"sort" schema:"sort" binding:"omitempty"`
	Dir    string `json:"dir" schema:"dir" binding:"omitempty"`
}

type FilterRatingSummary struct {
	SourceType string   `json:"source_type"`
	SourceUid  []string `json:"source_uid"`
	StoreUID   []string `json:"store_uid,omitempty"`
	RatingType []string `json:"rating_type"`
}

func (f FilterRatingSummary) ValidateSourceUID() *message.Message {
	if len(f.SourceUid) == 0 {
		return &message.ErrSourceUidRequire
	}

	if len(f.SourceUid) > 50 {
		return &message.ErrSourceUidRequire
	}

	return nil
}

func (f FilterRatingSummary) ValidateStoreUID() *message.Message {
	if len(f.StoreUID) == 0 {
		return &message.ErrStoreUidRequire
	}

	if len(f.StoreUID) > 20 {
		return &message.ErrStoreUidMax
	}

	return nil
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
	// Filter available {"user_id_legacy": [""], "source_trans_id": [""], "value": "", "is_with_media": true, "start_date": "", "end_date": "", ""}
	Filter string `json:"filter" schema:"filter" binding:"omitempty"`
	Limit  int    `json:"limit" schema:"limit" binding:"omitempty,numeric,min=1,max=100"`
	Page   int    `json:"page" schema:"page" binding:"omitempty,numeric,min=1"`
	Sort   string `json:"sort" schema:"sort" binding:"omitempty"`
	Dir    string `json:"dir" schema:"dir" binding:"omitempty"`
}

type FilterRatingSubmission struct {
	RatingID      []string     `json:"rating_id"`
	LikertFilter  LikertFilter `json:"likert_filter"`
	UserIdLegacy  []string     `json:"user_id_legacy"`
	SourceTransID []string     `json:"source_trans_id"`
	Value         string       `json:"value"`
	StartDate     string       `json:"start_date"`
	EndDate       string       `json:"end_date"`
}

func (req FilterRatingSubmission) ValidateFormatDate() error {
	return validation.ValidateStruct(&req,
		// Default
		validation.Field(&req.StartDate, validation.When(req.StartDate != "", validation.Date("2006-01-02").
			Error("start_date is invalid format, format should be 2006-01-02"))),
		validation.Field(&req.EndDate, validation.When(req.EndDate != "", validation.Date("2006-01-02").
			Error("end_date is invalid format, format should be 2006-01-02"))),
	)
}

type LikertFilter struct {
	RatingId string   `json:"rating_id"`
	Value    []string `json:"value"`
}

// swagger:parameters GetRatingBySourceTypeAndActor
type GetRatingBySourceTypeAndActorRequest struct {
	// in: path
	// required: true
	SourceType string `json:"source_type"`
	// in: path
	// required: true
	SourceUID string `json:"source_uid"`

	// Filter available {"rating_type": ["rating_like_dislike", "list_doctor_likert_for_positif_reviews", "list_doctor_likert_for_negative_reviews"]}
	Filter string `json:"filter" schema:"filter" binding:"omitempty"`
}

type GetRatingBySourceTypeAndActorFilter struct {
	RatingType []string `json:"rating_type"`
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
	interfaceAllSource := make([]interface{}, len(sourceType))
	for i, v := range sourceType {
		interfaceAllSource[i] = v
	}
	return validation.ValidateStruct(&req,
		validation.Field(&req.SourceType, validation.Required, validation.In(interfaceAllSource...).Error(fmt.Sprintf("source_type should be %s", strings.Join(sourceType, ",")))),
	)
}
