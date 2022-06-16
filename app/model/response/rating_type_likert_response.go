package response

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// swagger:model RatingTypeLikertResponse
type RatingTypeLikertResponse struct {
	ID            primitive.ObjectID `json:"id"`
	Type          string             `json:"type"`
	Description   string             `json:"description"`
	NumStatements int                `json:"num_statements"`
	Statement01   string             `json:"statement_01"`
	Statement02   string             `json:"statement_02"`
	Statement03   string             `json:"statement_03"`
	Statement04   string             `json:"statement_04"`
	Statement05   string             `json:"statement_05"`
	Statement06   string             `json:"statement_06"`
	Statement07   string             `json:"statement_07"`
	Statement08   string             `json:"statement_08"`
	Statement09   string             `json:"statement_09"`
	Statement10   string             `json:"statement_10"`
	Status        bool               `json:"status"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}
