package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:model RatingSubHelpfulCol
type RatingSubHelpfulCol struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	RatingSubmissionID string             `json:"rating_submission_id" bson:"rating_submission_id,omitempty"`
	UserID             string             `json:"user_id" bson:"user_id,omitempty"`
	UserIDLegacy       string             `json:"user_id_legacy" bson:"user_id_legacy,omitempty"`
	IPAddress          string             `json:"ip_address" bson:"ip_address,omitempty"`
	UserAgent          string             `json:"user_agent" bson:"user_agent,omitempty"`
	Status             bool               `json:"status" bson:"status,omitempty"`
	CreatedAt          time.Time          `json:"-" bson:"created_at,omitempty"`
	UpdatedAt          time.Time          `json:"-" bson:"updated_at,omitempty"`
}

func (RatingSubHelpfulCol) CollectionName() string {
	return "ratingSubHelpfulCol"
}
