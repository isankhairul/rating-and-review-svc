package request

// swagger:parameters getRatingFormulas
type GetRatingFormulasRequest struct {
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

type RatingFormulaFilter struct {
	TypeId     []string `json:"type_id"`
	SourceType []int    `json:"source_type"`
}

// swagger:parameters getRatingFormulaById deleteRatingFormula
type GetRatingFormulaRequest struct {
	// in: path
	// required: true
	Id string `json:"id"`
}

// swagger:parameters updateRatingFormula
type ReqUpdateRatingFormulaBody struct {
	// in: path
	// required: true
	Id string `json:"id"`
	//  in: body
	// required: true
	Body SaveRatingFormula `json:"body"`
}

// swagger:parameters createRatingFormula
type ReqCreateRatingFormulaBody struct {
	//  in: body
	// required: true
	Body SaveRatingFormula `json:"body"`
}

type SaveRatingFormula struct {
	// Source Type of rating formula
	// in: string
	SourceType string `json:"source_type,omitempty"`
	// Formula of rating formula
	// in: int
	Formula string `json:"formula,omitempty"`
	// Rating Type Id of rating formula
	// in: string
	RatingTypeId string `json:"rating_type_id,omitempty"`
	// Rating Type of rating formula
	// in: string
	RatingType string `json:"rating_type,omitempty"`
	// Status of rating formula
	// in: bool
	Status *bool `json:"status"`
	// For update
	Id string `json:"-"`
}
