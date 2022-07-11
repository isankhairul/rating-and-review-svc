package request

// swagger:parameters GetPublicListRatingSummaryRequest
type GetPublicListRatingSummaryRequest struct {
	// SourceType
	// in: path
	// required: true
	SourceType string `json:"source_type"`
	Filter     string `json:"filter" schema:"filter" binding:"omitempty"`
	Limit      int    `json:"limit" schema:"limit" binding:"omitempty,numeric,min=1,max=100"`
	Page       int    `json:"page" schema:"page" binding:"omitempty,numeric,min=1"`
	Sort       string `json:"sort" schema:"sort" binding:"omitempty"`
	Dir        string `json:"dir" schema:"dir" binding:"omitempty"`
}

type FilterRatingSummary struct {
	SourceUid  string   `json:"source_uid"`
	SourceType string   `json:"source_type"`
	RatingType []string `json:"rating_type"`
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
