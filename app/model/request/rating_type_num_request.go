package request

import (
	validation "github.com/itgelo/ozzo-validation/v4"
	"go-klikdokter/helper/message"
	"regexp"
)

var (
	regexType = "^[a-z0-9\\-._]+$"
)

// swagger:parameters SaveDoctorRequest
type ReqDoctorBody struct {
	//  in: body
	Body CreateRatingTypeNumRequest `json:"body"`
}

// swagger:parameters getRatingTypeNums
type GetRatingTypeNumsRequest struct {
	// Maximum records per page
	// in: query
	// type: integer
	Limit int64 ` json:"limit" schema:"limit" binding:"omitempty,numeric,min=1,max=100"`

	// Page No
	// in: query
	Page int ` json:"page" schema:"page" binding:"omitempty,numeric,min=1"`

	// Sort fields, example: name asc, uom desc
	// in: query
	Sort string ` json:"sort" schema:"sort" binding:"omitempty"`

	// Sort fields, example: name asc, uom desc
	// in: query
	Dir string ` json:"dir" schema:"dir" binding:"omitempty"`

	Filter string `json:"filter"`
}

type Filter struct {
	TypeId   []string `json:"type_id"`
	MinScore []int    `json:"min_score"`
	MaxScore []int    `json:"max_score"`
}

// swagger:parameters updateRatingTypeNum
type ReqUpdateRatingTypeNumberBody struct {
	// in: path
	// required: true
	Id string `json:"id"`
	//  in: body
	// required: true
	Body EditRatingTypeNumRequest `json:"body"`
}

// swagger:parameters createRatingTypeNum
type ReqCreateRatingTypeNumberBody struct {
	//  in: body
	// required: true
	Body CreateRatingTypeNumRequest `json:"body"`
}

type CreateRatingTypeNumRequest struct {
	// Type of rating type num
	// in: string
	Type string `json:"type"`
	// Description of rating type num
	// in: string
	Description *string `json:"description"`
	// Min Score of rating type num
	// in: integer
	MinScore *int `json:"min_score"`
	// Max Score of rating type num
	// in: integer
	MaxScore *int `json:"max_score"`
	// Scale of rating type num
	// in: integer
	Scale *int `json:"scale"`
	// Intervals of rating type num
	// in: integer
	Intervals int `json:"intervals"`
	// Status of rating type num
	// in: bool
	Status *bool `json:"status"`
	// for update
	Id string `json:"-"`
}

type EditRatingTypeNumRequest struct {
	// Type of rating type num
	// in: string
	Type string `json:"type"`
	// Description of rating type num
	// in: string
	Description *string `json:"description"`
	// Min Score of rating type num
	// in: integer
	MinScore *int `json:"min_score"`
	// Max Score of rating type num
	// in: integer
	MaxScore *int `json:"max_score"`
	// Scale of rating type num
	// in: integer
	Scale *int `json:"scale"`
	// Intervals of rating type num
	// in: integer
	Intervals *int `json:"intervals"`
	// Status of rating type num
	// in: bool
	Status *bool `json:"status"`
	// for update
	Id string `json:"-"`
}

// swagger:parameters getRatingTypeNumById deleteRatingTypeNum
type GetRatingTypeNumRequest struct {
	// in: path
	// required: true
	Id string `json:"id"`
}

func (req CreateRatingTypeNumRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Type, validation.Required.Error(message.ErrTypeReq.Message), validation.Match(regexp.MustCompile(regexType)).Error(message.ErrTypeFormatReq.Message)),
		validation.Field(&req.MinScore, validation.NotNil.Error((message.ErrMinScoreReq.Message))),
		validation.Field(&req.MaxScore, validation.NotNil.Error((message.ErrMaxScoreReq.Message))),
		validation.Field(&req.Scale, validation.NotNil.Error((message.ErrScaleReq.Message))),
	)
}
