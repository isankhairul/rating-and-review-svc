package request

// swagger:parameters GetRatingBySourceTypeAndActor
type GetRatingBySourceTypeAndActorRequest struct {
	// in: path
	// required: true
	SourceType string `json:"source_type"`
	// in: path
	// required: true
	SourceUID string `json:"source_uid"`
}
