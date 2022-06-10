package encoder

import (
	"context"
	"encoding/json"
	"go-klikdokter/app/model/base"
	"go-klikdokter/helper/message"
	"net/http"
)

type errorer interface {
	error() error
}

type Empty struct{}

type errorResponse struct {
	// Meta is the API response information
	// in: struct{}
	Meta struct {
		// Code is the response code
		//in: int
		Code int `json:"code"`
		// Message is the response message
		//in: string
		Message string `json:"message"`
	} `json:"meta"`
	// Data is our data
	// in: struct{}
	Data interface{} `json:"data"`
	// Errors is the response message
	//in: string
	Errors interface{} `json:"errors,omitempty"`
}

func EncodeResponseHTTP(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	if err, ok := resp.(errorer); ok && err.error() != nil {
		EncodeError(ctx, err.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	result := base.GetHttpResponse(resp)
	code := result.Meta.Code
	switch code {
	case message.DataNotFoundCode:
		w.WriteHeader(http.StatusNotFound)
	case message.ErrNoAuth.Code:
		w.WriteHeader(http.StatusUnauthorized)
	case message.BadRequestCode, message.DataNotFoundCode, message.ErrSaveData.Code, message.ErrTypeReq.Code,
		message.ErrTypeFormatReq.Code, message.ErrIdFormatReq.Code, message.ErrScaleValueReq.Code,
		message.ErrIntervalsValueReq.Code, message.UserAgentTooLong.Code:
		w.WriteHeader(http.StatusBadRequest)
	case message.SuccessCode, message.ErrNoData.Code:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	return json.NewEncoder(w).Encode(resp)
}

//Encode error, for HTTP
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	result := &errorResponse{}
	result.Meta.Code = message.ErrReq.Code
	result.Meta.Message = err.Error()
	result.Data = Empty{}
	_ = json.NewEncoder(w).Encode(result)
}
