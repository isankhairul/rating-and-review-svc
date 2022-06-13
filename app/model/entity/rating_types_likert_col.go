package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// swagger:model RatingTypesLikertCol
type RatingTypesLikertCol struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Type          string             `json:"type" bson:"type,omitempty"`
	Description   string             `json:"description" bson:"description,omitempty"`
	NumStatements int                `json:"num_statements" bson:"num_statements,omitempty"`
	Statement01   *string            `json:"statement_01" bson:"statement_01"`
	Statement02   *string            `json:"statement_02" bson:"statement_02"`
	Statement03   *string            `json:"statement_03" bson:"statement_03"`
	Statement04   *string            `json:"statement_04" bson:"statement_04"`
	Statement05   *string            `json:"statement_05" bson:"statement_05"`
	Statement06   *string            `json:"statement_06" bson:"statement_06"`
	Statement07   *string            `json:"statement_07" bson:"statement_07"`
	Statement08   *string            `json:"statement_08" bson:"statement_08"`
	Statement09   *string            `json:"statement_09" bson:"statement_09"`
	Statement10   *string            `json:"statement_10" bson:"statement_10"`
	Status        *bool              `json:"status" bson:"status,omitempty"`
	CreatedAt     *time.Time         `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt     *time.Time         `json:"updated_at" bson:"updated_at,omitempty"`
}
