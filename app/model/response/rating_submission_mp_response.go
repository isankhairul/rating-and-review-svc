package response

import "go.mongodb.org/mongo-driver/bson/primitive"

type RatingSubmissionMpResponse struct {
	RatingID      string             `json:"rating_id" bson:"rating_id"`
	UserID        *string            `json:"user_id" bson:"user_id"`
	UserIDLegacy  *string            `json:"user_id_legacy" bson:"user_id_legacy"`
	Comment       string             `json:"comment" bson:"comment"`
	Value         string             `json:"value" bson:"value"`
	SourceTransID string             `json:"source_trans_id" bson:"source_trans_id"`
	Media         []MediaObjResponse `json:"media" bson:"media"`
	IsWithMedia   bool               `json:"is_with_media" bson:"is_with_media"`
}

type MediaObjResponse struct {
	UID        string `json:"uid"`
	MediaPath  string `json:"media_path"`
	MediaImage string `json:"media_image"`
}

type CreateRatingSubmissionMpResponse struct {
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
