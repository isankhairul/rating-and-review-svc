package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RatingSubmissionMp struct {
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
	SourceUID       string             `json:"source_uid" bson:"source_uid"`
	SourceType      string             `json:"source_type" bson:"source_type"`
	SourceTransID   string             `json:"source_trans_id" bson:"source_trans_id,omitempty"`
	LikeCounter     int                `json:"like_counter" bson:"like_counter"`
	UserPlatform    string             `json:"user_platform" bson:"user_platform,omitempty"`
	Cancelled       bool               `json:"cancelled" bson:"cancelled"`
	CancelledReason string             `json:"cancelled_reason" bson:"cancelled_reason"`
	IsAnonymous     bool               `json:"is_anonymous" bson:"is_anonymous,omitempty"`
	Media           []MediaObj         `json:"media" bson:"media"`
	IsWithMedia     bool               `json:"is_with_media" bson:"is_with_media"`
	OrderNumber     string             `json:"order_number" bson:"order_number"`
	CreatedAt       time.Time          `json:"-" bson:"created_at,omitempty"`
	UpdatedAt       time.Time          `json:"-" bson:"updated_at,omitempty"`
	RatingTypeID    string             `json:"rating_type_id" bson:"rating_type_id, omitempty"`
	Reply           string             `json:"reply" bson:"reply"`
	ReplyBy         string             `json:"reply_by" bson:"reply_by"`
}

func (RatingSubmissionMp) CollectionName() string {
	return "ratingSubMpCol"
}

type MediaObj struct {
	UID       string `json:"uid" bson:"uid"`
	MediaPath string `json:"media_path" bson:"media_path"`
}
