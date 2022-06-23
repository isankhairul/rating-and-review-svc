package response

type RatingSummaryResponse struct {
	SourceUID   string  `json:"source_uid" bson:"source_uid"`
	TotalReview int     `json:"total_review" bson:"total_review"`
	Value       float64 `json:"value" bson:"value"`
}
