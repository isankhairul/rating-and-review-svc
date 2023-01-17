package transport

import (
	"context"
	"go-klikdokter/app/api/endpoint"
	"go-klikdokter/app/model/base/encoder"
	request_image "go-klikdokter/app/model/request/upload"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/_struct"
	"go-klikdokter/pkg/util"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"go-klikdokter/app/middleware"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func UploadHttpHandler(s service.UploadService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeUploadImageEndpoint(s, logger)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	pr.Methods(http.MethodPost).Path(_struct.PrefixBase + "/upload/images").Handler(httptransport.NewServer(
		ep.MakeUploadImage,
		decodeUploadImage,
		encoder.EncodeResponseHTTPWithCorrelationID,
		append(options, httptransport.ServerBefore(middleware.CorrelationIdToContext()))...,
	))

	return pr
}

func decodeUploadImage(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request_image.UploadImageRequest
	file, handler, err := r.FormFile("image")
	if err != nil {
		return nil, err
	}

	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)

	baseValidation := util.NewValidationImage(handler.Filename, fileBytes, &util.DefaultMimeAllowed, &util.DefaultSizeAllowed)

	if err := baseValidation.ValidateSizeAndMime(); err != nil {
		return nil, err
	}
	mimetype := http.DetectContentType(fileBytes)
	filename := strconv.FormatInt(time.Now().Unix(), 10) + "-" + handler.Filename
	req.FileName = filename
	req.Image = fileBytes
	req.SourceType = r.FormValue("source_type")
	req.SourceUid = r.FormValue("source_uid")
	req.MimeType = mimetype
	return req, nil
}
