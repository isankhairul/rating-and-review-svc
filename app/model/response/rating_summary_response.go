package response

type RatingSummaryResponse struct {
	RatingID     string  `json:"rating_id" bson:"rating_id"`
	UserID       *string `json:"user_id" bson:"user_id"`
	UserIDLegacy *string `json:"user_id_legacy" bson:"user_id_legacy"`
	SourceUID    string  `json:"source_uid" bson:"source_uid"`
	Name         string  `json:"name" bson:"name"`
	TotalReview  int     `json:"total_review" bson:"total_review"`
	Value        float64 `json:"value" bson:"value"`
}
