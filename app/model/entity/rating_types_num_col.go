package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// swagger:model RatingTypesNumCol
type RatingTypesNumCol struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Type        string             `json:"type,omitempty" bson:"type,omitempty"`
	Description *string            `json:"description,omitempty" bson:"description,omitempty"`
	MinScore    *int               `json:"min_score,omitempty" bson:"min_score,omitempty"`
	MaxScore    *int               `json:"max_score,omitempty" bson:"max_score,omitempty"`
	Scale       *int               `json:"scale,omitempty" bson:"scale,omitempty"`
	Intervals   *int               `json:"intervals,omitempty" bson:"intervals,omitempty"`
	Status      *bool              `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}

func (RatingTypesNumCol) CollectionName() string {
	return "ratingTypesNumCol"
}
