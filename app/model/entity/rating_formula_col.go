package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:model RatingFormulaCol
type RatingFormulaCol struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SourceType   string             `json:"source_type" bson:"source_type,omitempty"`
	Formula      string             `json:"formula" bson:"formula,omitempty"`
	RatingTypeId string             `json:"rating_type_id" bson:"rating_type_id,omitempty"`
	RatingType   string             `json:"rating_type" bson:"rating_type,omitempty"`
	Status       *bool              `json:"status" bson:"status,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}
