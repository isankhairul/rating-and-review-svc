package response

import "go.mongodb.org/mongo-driver/bson/primitive"

type RatingSubmissionMpResponse struct {
	RatingID      string   `json:"rating_id" bson:"rating_id"`
	UserID        *string  `json:"user_id" bson:"user_id"`
	UserIDLegacy  *string  `json:"user_id_legacy" bson:"user_id_legacy"`
	Comment       string   `json:"comment" bson:"comment"`
	Value         string   `json:"value" bson:"value"`
	SourceTransID string   `json:"source_trans_id" bson:"source_trans_id"`
	MediaPath     []string `json:"media_path" bson:"media_path"`
	IsWithMedia   bool     `json:"is_with_media" bson:"is_with_media"`
}

type CreateRatingSubmissionMpResponse struct {
	ID                primitive.ObjectID `json:"id"`
	RatingID          string             `json:"rating_id,omitempty"`
	RatingDescription string             `json:"rating_decription,omitempty"`
	Value             string             `json:"value,omitempty"`
	// MediaPath         []string           `json:"media_path"`
	IsWithMedia       bool               `json:"is_with_media"`
	OrderNumber       string             `json:"order_number"`
}

type RatingSummaryMpResponse struct {
	ID            primitive.ObjectID `json:"id"`
	Name          string             `json:"name,omitempty"`
	Description   *string            `json:"description,omitempty"`
	SourceUid     string             `json:"source_uid,omitempty"`
	SourceType    string             `json:"source_type,omitempty"`
	RatingType    string             `json:"rating_type,omitempty"`
	RatingTypeId  string             `json:"rating_type_id,omitempty"`
	RatingSummary interface{}        `json:"rating_summary,omitempty"`
}
