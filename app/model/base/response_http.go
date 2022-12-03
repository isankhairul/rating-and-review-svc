package base

import (
	"context"
	"fmt"
	"reflect"
)

// swagger:model SuccessResponse
type responseHttp struct {
	// Meta is the API response information
	// in: MetaResponse
	Meta metaResponse `json:"meta"`
	// Pagination of the paginate respons
	// in: PaginationResponse
	Pagination *Pagination `json:"pagination,omitempty"`
	// Data is our data
	// in: DataResponse
	Data data `json:"data"`
	// Errors is the response message
	//in: string
	Errors interface{} `json:"errors,omitempty"`
}

// swagger:model SuccessResponse
type responseHttpWithCtx struct {
	// Meta is With API response information
	// in: MetaResponse
	Meta metaResponseWithCtx `json:"meta"`
	// Pagination of the paginate respons
	// in: PaginationResponse
	Pagination *Pagination `json:"pagination,omitempty"`
	// Data is our data
	// in: DataResponse
	Data data `json:"data"`
	// Errors is the response message
	//in: string
	Errors interface{} `json:"errors,omitempty"`
}

// swagger:model MetaResponse
type metaResponse struct {
	// Code is the response code
	//in: int
	Code int `json:"code"`
	// Message is the response message
	//in: string
	Message string `json:"message"`
}

// swagger:model MetaResponseCtx
type metaResponseWithCtx struct {
	// CorrelationId is the response correlation_id
	//in: string
	CorrelationId string `json:"correlation_id"`
	// Code is the response code
	//in: int
	Code int `json:"code"`
	// Message is the response message
	//in: string
	Message string `json:"message"`
}

// swagger:model DataResponse
type data struct {
	Records interface{} `json:"records,omitempty"`
	Record  interface{} `json:"record,omitempty"`
}

func SetHttpResponse(code int, message string, result interface{}, paging *Pagination) interface{} {
	dt := data{}
	isSlice := reflect.ValueOf(result).Kind() == reflect.Slice
	if isSlice {
		dt.Records = result
		dt.Record = nil
	} else {
		dt.Records = nil
		dt.Record = result
	}

	return responseHttp{
		Meta: metaResponse{
			Code:    code,
			Message: message,
		},
		Pagination: paging,
		Data:       dt,
	}
}

func SetHttpResponseWithCtx(ctx context.Context,code int, message string, result interface{}, paging *Pagination, errMsg interface{}) interface{} {
	dt := data{}
	isSlice := reflect.ValueOf(result).Kind() == reflect.Slice
	if isSlice {
		dt.Records = result
		dt.Record = nil
	} else {
		dt.Records = nil
		dt.Record = result
	}

	return responseHttpWithCtx{
		Meta: metaResponseWithCtx{
			CorrelationId: fmt.Sprint(ctx.Value(CorrelationIdContextKey)),
			Code:    code,
			Message: message,
		},
		Pagination: paging,
		Data:       dt,
	}
}

func GetHttpResponse(resp interface{}) *responseHttp {
	result, ok := resp.(responseHttp)

	if ok {
		return &result
	}
	return nil
}
