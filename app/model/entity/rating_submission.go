package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RatingSubmisson struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	RatingID        string             `json:"rating_id" bson:"rating_id,omitempty"`
	UserID          *string            `json:"user_id" bson:"user_id,omitempty"`
	UserIDLegacy    *string            `json:"user_id_legacy" bson:"user_id_legacy,omitempty"`
	DisplayName     *string            `json:"display_name" bson:"display_name,omitempty"`
	Comment         *string            `json:"comment" bson:"comment,omitempty"`
	Value           string             `json:"value" bson:"value,omitempty"`
	IPAddress       string             `json:"ip_address" bson:"ip_address,omitempty"`
	UserAgent       string             `json:"user_agent" bson:"user_agent,omitempty"`
	Avatar          string             `json:"avatar" bson:"avatar,omitempty"`
	SourceTransID   string             `json:"source_trans_id" bson:"source_trans_id,omitempty"`
	LikeCounter     int                `json:"like_counter" bson:"like_counter,omitempty"`
	UserPlatform    string             `json:"user_platform" bson:"user_platform,omitempty"`
	Cancelled       bool               `json:"cancelled" bson:"cancelled,omitempty"`
	CancelledReason string             `json:"cancelled_reason" bson:"cancelled_reason,omitempty"`
	IsAnonymous     bool               `json:"is_anonymous" bson:"is_anonymous,omitempty"`
	MediaPath       []string           `json:"media_path" bson:"media_path,omitempty"`
	IsWithMedia     bool               `json:"is_with_media" bson:"is_with_media,omitempty"`
	CreatedAt       time.Time          `json:"-" bson:"created_at,omitempty"`
	UpdatedAt       time.Time          `json:"-" bson:"updated_at,omitempty"`
}

func (RatingSubmisson) CollectionName() string {
	return "ratingSubCol"
}
