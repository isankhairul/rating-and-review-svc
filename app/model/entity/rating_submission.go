package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type RatingSubmisson struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	RatingID     string             `json:"rating_id" bson:"rating_id"`
	UserID       *string            `json:"user_id" bson:"user_id"`
	UserIDLegacy *string            `json:"user_id_legacy" bson:"user_id_legacy"`
	Comment      string             `json:"comment" bson:"comment"`
	Value        float64            `json:"value" bson:"value"`
	IPAddress    string             `json:"ip_address" bson:"ip_address"`
	UserAgent    string             `json:"user_agent" bson:"user_agent"`
	CreatedAt    time.Time          `json:"-" bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `json:"-" bson:"updated_at,omitempty"`
}
