package util_media

import (
	"encoding/json"
	"fmt"
	"go-klikdokter/app/model/request"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/httputil"
	"net/http"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/viper"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ImageHouseKeeping(logg log.Logger, mediaPath []request.MediaPathObj, ratingSubsID string) ([]map[string]interface{}, error) {
	logger := log.With(logg, "media-svc", "ImageHouseKeeping")
	response := []map[string]interface{}{}
	var error error
	token, _ := global.GenerateJwt()

	if len(mediaPath) > 0 {
		for _, mp := range mediaPath {
			param := map[string]string{
				"source_type": "rnr",
				"source_uid":  ratingSubsID,
			}
			jsonData, _ := json.Marshal(param)

			url := viper.GetString("media-service.url-image-house-keeping") + "/" + mp.UID
			httpRequest := httputil.NewHttpRequest()
			headers := map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			}
			statusCode, _, dataResponse, err := httpRequest.PerformRequest(http.MethodPut, url, jsonData, nil, headers, logg)
			jsonDataResponse, _ := json.Marshal(dataResponse)

			response = append(response, dataResponse)

			if statusCode != 200 || err != nil {
				_ = level.Error(logger).Log("Log", fmt.Sprintf("Got error from endpoint: %s, with request: %v, response: %v",
					url, string(jsonData), string(jsonDataResponse)))
				error = err
				continue
			}
			_ = level.Info(logger).Log(fmt.Sprintf("Success Push, with Response: %v", err))
		}
	}

	return response, error
}
