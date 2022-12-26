package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-klikdokter/app/middleware"
	upload_request "go-klikdokter/app/model/request/upload"
	"go-klikdokter/app/model/response"
	global "go-klikdokter/helper/global"
	"go-klikdokter/helper/httphelper"
	"go-klikdokter/helper/message"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	resty "github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type UploadService interface {
	UploadImage(ctx context.Context ,input upload_request.UploadImageRequest) (response.UploadResponse, message.Message, interface{})
}

type uploadServiceImpl struct {
	logger 		log.Logger
}

func NewUploadService(
	lg log.Logger, 
) UploadService {
	return &uploadServiceImpl{lg}
}

// swagger:operation POST /upload/images Image ReqUploadForm
//
// Upload Image Rating
//
// This will create a new image
//
// ---
// security:
// - Bearer: []
// consumes: [multipart/form-data]
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//          $ref: '#/definitions/metaResponseWithCorrelationID'
//         data:
//           properties:
//             records:
//               $ref: '#/definitions/UploadResponse'
//           type: object
func (updSvc *uploadServiceImpl) UploadImage(ctx context.Context, input upload_request.UploadImageRequest) (response.UploadResponse, message.Message, interface{}) {
	result  := response.UploadResponse{}
	var messages message.Message
	errMsg := make(map[string]interface{})

	client := resty.New()
    
	correlationId := fmt.Sprint(ctx.Value(middleware.CorrelationIdContextKey))
	logger := log.With(updSvc.logger, "MediaService", fmt.Sprint("performRequest-", correlationId))
	token, _ := global.GenerateJwt()

	urlMediaService := viper.GetString("media-service.url")
	mapFormData := map[string]string{
		"name" : input.FileName,
		"source_type" : "rnr",
		"media_category_uid" : viper.GetString("media-service.media-category-uid"),
		"description" : "direct upload from rnr",
	}
	resp, err := client.R().
		SetHeader("Authorization", "Bearer " + token).
		SetHeader("X-Correlation-ID", correlationId).
		SetFileReader("image", input.FileName, bytes.NewReader(input.Image)).
		SetFormData(mapFormData).
		Post(urlMediaService)
	level.Info(logger).Log("type","[Media-Svc]", "respStatus", resp.StatusCode() ,"respBody", string(resp.Body()))
	if err != nil {
		errMsg["image"] = "Failed to Upload Image"
		return result, message.ErrUploadMedia, nil	
	}
    
	rspStatusCode := resp.StatusCode()

    var ResponseMedia response.ResponseHttpMedia
	json.Unmarshal(resp.Body(), &ResponseMedia)
	
	if rspStatusCode == 200 {
		result.UID = &ResponseMedia.Data.Record.Uid
		result.MediaPath = &ResponseMedia.Data.Record.ImageFiles[0].MediaPath
		messages = message.SuccessMsg
	} else {
		result.UID = nil
		result.MediaPath= nil
		errMsg["image"] = "Failed to Upload Image"
		messages = message.ErrUploadMedia
	}
	return result, messages, errMsg
}

func (updSvc *uploadServiceImpl) UploadImageOld(ctx context.Context, input upload_request.UploadImageRequest) (response.UploadResponse, message.Message, interface{}) {
	result  := response.UploadResponse{}
	var messages message.Message
	errMsg := make(map[string]interface{})
    //new  multipart writer
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)

    fw, _ := writer.CreateFormField("name")
    io.Copy(fw, strings.NewReader(input.FileName))

	fw, _ = writer.CreateFormField("source_type")
    io.Copy(fw, strings.NewReader("rnr"))

	mediaCategoryUid := viper.GetString("media-service.media-category-uid")
    fw, _ = writer.CreateFormField("media_category_uid")
    io.Copy(fw, strings.NewReader(mediaCategoryUid))

    fileReader := bytes.NewReader(input.Image)    
    mimetype := input.MimeType
	fw, _ = CreateCustomFormFile(writer, input.FileName, mimetype)
	io.Copy(fw, fileReader)

    writer.Close()

	correlationId := fmt.Sprint(ctx.Value(middleware.CorrelationIdContextKey))
	logger := log.With(updSvc.logger, "MediaService", fmt.Sprint("performRequest-", correlationId))
	var queryParams map[string] string
	token, _ := global.GenerateJwt()
	headers := map[string]string{
		"X-Correlation-ID": correlationId,
		"Authorization":    fmt.Sprint("Bearer ", token),
	}

	urlMediaService := viper.GetString("media-service.url")
	rspStatusCode, data, err := httphelper.PerformRequestMultipartWithLog(logger, "POST", urlMediaService, body.Bytes(), queryParams, headers, writer)

	if err != nil {
		errMsg["image"] = "Failed to Upload Image"
		return result, message.ErrUploadMedia, nil	
	}
    
    var ResponseMedia response.ResponseHttpMedia
	json.Unmarshal(data, &ResponseMedia)
	
	if rspStatusCode == 200 {
		result.UID = &ResponseMedia.Data.Record.Uid
		result.MediaPath = &ResponseMedia.Data.Record.ImageFiles[0].MediaPath
		messages = message.SuccessMsg
	} else {
		result.UID = nil
		result.MediaPath= nil
		errMsg["image"] = "Failed to Upload Image"
		messages = message.ErrUploadMedia
	}
	return result, messages, errMsg
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func CreateCustomFormFile(w *multipart.Writer, filename string, mimeType string) (io.Writer, error) {
    h := make(textproto.MIMEHeader)
    h.Set("Content-Disposition",
        fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
        "image", escapeQuotes(filename)))
    h.Set("Content-Type", mimeType)
    return w.CreatePart(h)
}

func RemoveLocalImage(filename string) {
	// Remove local file
	err := os.Remove(filename)
	fmt.Println("LOI XOA FILE: ", err)
}