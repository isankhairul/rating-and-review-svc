package response

import (
	"go-klikdokter/app/model/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:model RatingTypeLikertResponse
type RatingBySourceTypeAndActorResponse struct {
	SourceUID  string        `json:"source_uid"`
	SourceType string        `json:"source_type"`
	Ratings    []interface{} `json:"ratings"`
}

type PublicRatingNumericResponse struct {
	Type              string             `json:"type"`
	RatingId          primitive.ObjectID `json:"rating_id"`
	RatingType        string             `json:"rating_type"`
	RatingDescription string             `json:"rating_description,omitempty"`
	RatingMinScore    int                `json:"rating_min_score"`
	RatingMaxScore    int                `json:"rating_max_score"`
}

type PublicRatingLikertResponse struct {
	Type                string             `json:"type"`
	RatingId            primitive.ObjectID `json:"rating_id"`
	RatingType          string             `json:"rating_type"`
	RatingDescription   string             `json:"rating_description"`
	RatingNumStatements int                `json:"rating_num_statements"`
	RatingStatement01   *string            `json:"rating_statement_01,omitempty"`
	RatingStatement02   *string            `json:"rating_statement_02,omitempty"`
	RatingStatement03   *string            `json:"rating_statement_03,omitempty"`
	RatingStatement04   *string            `json:"rating_statement_04,omitempty"`
	RatingStatement05   *string            `json:"rating_statement_05,omitempty"`
	RatingStatement06   *string            `json:"rating_statement_06,omitempty"`
	RatingStatement07   *string            `json:"rating_statement_07,omitempty"`
	RatingStatement08   *string            `json:"rating_statement_08,omitempty"`
	RatingStatement09   *string            `json:"rating_statement_09,omitempty"`
	RatingStatement10   *string            `json:"rating_statement_10,omitempty"`
}

func MapRatingNumericToRatingNumericResp(data entity.RatingTypesNumCol) *PublicRatingNumericResponse {
	return &PublicRatingNumericResponse{
		Type:              "numeric",
		RatingId:          data.ID,
		RatingType:        data.Type,
		RatingDescription: *data.Description,
		RatingMinScore:    *data.MinScore,
		RatingMaxScore:    *data.MaxScore,
	}
}

func MapRatingLikertToRatingNumericResp(data entity.RatingTypesLikertCol) *PublicRatingLikertResponse {
	result := PublicRatingLikertResponse{}
	result.Type = "likert"
	result.RatingId = data.ID
	result.RatingType = data.Type
	result.RatingDescription = *data.Description
	result.RatingNumStatements = data.NumStatements
	if data.Statement01 != nil && len(*data.Statement01) != 0 {
		result.RatingStatement01 = data.Statement01
	} else {
		result.RatingStatement01 = nil
	}

	if data.Statement02 != nil && len(*data.Statement02) != 0 {
		result.RatingStatement02 = data.Statement02
	} else {
		result.RatingStatement02 = nil
	}

	if data.Statement03 != nil && len(*data.Statement03) != 0 {
		result.RatingStatement03 = data.Statement03
	} else {
		result.RatingStatement03 = nil
	}

	if data.Statement04 != nil && len(*data.Statement04) != 0 {
		result.RatingStatement04 = data.Statement04
	} else {
		result.RatingStatement04 = nil
	}

	if data.Statement05 != nil && len(*data.Statement05) != 0 {
		result.RatingStatement05 = data.Statement05
	} else {
		result.RatingStatement05 = nil
	}

	if data.Statement06 != nil && len(*data.Statement06) != 0 {
		result.RatingStatement06 = data.Statement06
	} else {
		result.RatingStatement06 = nil
	}

	if data.Statement07 != nil && len(*data.Statement07) != 0 {
		result.RatingStatement07 = data.Statement07
	} else {
		result.RatingStatement07 = nil
	}

	if data.Statement08 != nil && len(*data.Statement08) != 0 {
		result.RatingStatement08 = data.Statement08
	} else {
		result.RatingStatement08 = nil
	}

	if data.Statement09 != nil && len(*data.Statement09) != 0 {
		result.RatingStatement09 = data.Statement09
	} else {
		result.RatingStatement09 = nil
	}

	if data.Statement10 != nil && len(*data.Statement10) != 0 {
		result.RatingStatement10 = data.Statement10
	} else {
		result.RatingStatement10 = nil
	}
	return &result
}
