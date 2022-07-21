package request

import (
	"go-klikdokter/helper/message"
	"regexp"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters updateRatingTypeLikert
type ReqUpdateRatingTypeLikertBody struct {
	// in: path
	// required: true
	Id string `json:"id"`
	//  in: body
	// required: true
	Body SaveRatingTypeLikertRequest `json:"body"`
}

// swagger:parameters createRatingTypeLikertRequest
type ReqCreateRatingTypeLikertBody struct {
	//  in: body
	// required: true
	Body SaveRatingTypeLikertRequest `json:"body"`
}

type SaveRatingTypeLikertRequest struct {
	// Type of rating type likert
	// in: string
	Type string `json:"type,omitempty"`
	// Description of rating type likert
	// in: string
	Description *string `json:"description,omitempty"`
	// NumStatements of rating type likert
	// in: integer
	NumStatements int `json:"num_statements,omitempty"`
	// Statement 1 of rating type likert
	// in: string
	Statement01 *string `json:"statement_01,omitempty"`
	// Statement 2 of rating type likert
	// in: string
	Statement02 *string `json:"statement_02,omitempty"`
	// Statement 3 of rating type likert
	// in: string
	Statement03 *string `json:"statement_03,omitempty"`
	// Statement 4 of rating type likert
	// in: string
	Statement04 *string `json:"statement_04,omitempty"`
	// Statement 5 of rating type likert
	// in: string
	Statement05 *string `json:"statement_05,omitempty"`
	// Statement 6 of rating type likert
	// in: string
	Statement06 *string `json:"statement_06,omitempty"`
	// Statement 7 of rating type likert
	// in: string
	Statement07 *string `json:"statement_07,omitempty"`
	// Statement 8 of rating type likert
	// in: string
	Statement08 *string `json:"statement_08,omitempty"`
	// Statement 9 of rating type likert
	// in: string
	Statement09 *string `json:"statement_09,omitempty"`
	// Statement 10 of rating type likert
	// in: string
	Statement10 *string `json:"statement_10,omitempty"`
	// Status of rating type likert
	// in: bool
	Status *bool `json:"status"`
	// For update
	Id string `json:"-"`
}

// swagger:parameters getRatingTypeLikerts
type GetRatingTypeLikertsRequest struct {
	// Maximum records per page
	// in: query
	// type: integer
	Limit int64 ` json:"limit" schema:"limit" binding:"omitempty,numeric,min=1,max=100"`

	// Page No
	// in: query
	Page int ` json:"page" schema:"page" binding:"omitempty,numeric,min=1"`

	// Sort fields
	// in: query
	Sort string ` json:"sort" schema:"sort" binding:"omitempty"`

	// Sort direction asc or desc
	// in: query
	Dir string ` json:"dir" schema:"dir" binding:"omitempty"`

	Filter string `json:"filter"`
}

type FilterRatingTypeLikert struct {
	TypeId []string `json:"type_id"`
}

// swagger:parameters getRatingTypeLikertById deleteRatingTypeLikert
type GetRatingTypeLikertRequest struct {
	// in: path
	// required: true
	Id string `json:"id"`
}

func (req SaveRatingTypeLikertRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Type, validation.Match(regexp.MustCompile(regexType)).Error(message.ErrTypeFormatReq.Message)),
		validation.Field(&req.Type, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.NumStatements, validation.Required.Error(message.ErrReq.Message)),
	)
}
