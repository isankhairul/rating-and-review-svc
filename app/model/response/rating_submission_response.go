package response

type RatingSubmissonResponse struct {
	RatingID     string  `json:"rating_id" bson:"rating_id"`
	UserID       *string `json:"user_id" bson:"user_id"`
	UserIDLegacy *string `json:"user_id_legacy" bson:"user_id_legacy"`
	Comment      string  `json:"comment" bson:"comment"`
	Value        string  `json:"value" bson:"value"`
	SourTransID  string  `json:"source_trans_id" bson:"source_trans_id"`
}
