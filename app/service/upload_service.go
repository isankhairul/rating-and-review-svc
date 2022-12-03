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
	"github.com/spf13/viper"
)

type UploadService interface {
	UploadImage(ctx context.Context ,input upload_request.UploadImageRequest) (*response.UploadResponse, message.Message, interface{})
}

type uploadServiceImpl struct {
	logger 		log.Logger
}

func NewUploadService(
	lg log.Logger, 
) UploadService {
	return &uploadServiceImpl{lg}
}

// swagger:operation POST /updalo/images Image ReqUploadForm
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
//          $ref: '#/definitions/MetaResponseWithCtx'
//         data:
//           properties:
//             records:
//               $ref: '#/definitions/UploadResponse'
//           type: object
func (updSvc *uploadServiceImpl) UploadImage(ctx context.Context, input upload_request.UploadImageRequest) (*response.UploadResponse, message.Message, interface{}) {
	var result  *response.UploadResponse
    //new  multipart writer.
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    fw, _ := writer.CreateFormField("name")
    io.Copy(fw, strings.NewReader(input.FileName))
	mediaCategoryUid := viper.GetString("media-service.media-category-uid")
    fw, _ = writer.CreateFormField("media_category_uid")
    io.Copy(fw, strings.NewReader(mediaCategoryUid))

    fileReader := bytes.NewReader(input.Image)
    /**/
    // fileReader = bytes.NewReader(files)
    mimetype := input.MimeType
	fw, _ = CreateCustomFormFile(writer, input.FileName, mimetype)
	io.Copy(fw, fileReader)
    /**/

    writer.Close()
	correlationId := fmt.Sprint(ctx.Value(middleware.CorrelationIdContextKey))
	logger := log.With(updSvc.logger, "MediaService", fmt.Sprint("performRequest-", correlationId))
	var queryParams map[string] string
	token, _ := global.GenerateJwt()
	fmt.Println("token", token)
	headers := map[string]string{
		"X-Correlation-ID": correlationId,
		"Authorization":    fmt.Sprint("Bearer ", token),
	}
	urlMediaService := viper.GetString("media-service.url")
	rspStatusCode, data, _ := httphelper.PerformRequestMultipartWithLog(logger, "POST", urlMediaService, body.Bytes(), queryParams, headers, writer)
    //bodyString := string(data)
    var ResponseMedia response.ResponseHttpMedia
	json.Unmarshal(data, &ResponseMedia)
	if rspStatusCode == 200 {
		result.UID = ResponseMedia.Data.Record.Uid
		result.MediaPath= ResponseMedia.Data.Record.ImageFiles[0].MediaPath
	}
	return nil, message.SuccessMsg, nil
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