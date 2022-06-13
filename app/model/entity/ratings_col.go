package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// swagger:model Rating
type RatingsCol struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name           string             `json:"name,omitempty" bson:"name,omitempty"`
	Description    string             `json:"description,omitempty" bson:"description,omitempty"`
	SourceUid      string             `json:"source_uid,omitempty" bson:"source_uid,omitempty"`
	SourceType     string             `json:"source_type,omitempty" bson:"source_type,omitempty"`
	RatingType     string             `json:"rating_type,omitempty" bson:"rating_type,omitempty"`
	RatingTypeId   string             `json:"rating_type_id,omitempty" bson:"rating_type_id,omitempty"`
	CommentAllowed *bool              `json:"comment_allowed,omitempty" bson:"comment_allowed,omitempty"`
	Status         *bool              `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt      time.Time          `json:"-" bson:"created_at,omitempty"`
	UpdatedAt      time.Time          `json:"-" bson:"updated_at,omitempty"`
}

func (RatingsCol) CollectionName() string {
	return "ratingsCol"
}
