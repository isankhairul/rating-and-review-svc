package thumbor

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-klikdokter/helper/config"
	"net/url"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/datatypes"
)

func GetThumborUrl(input string) string {
	baseUrl := config.GetConfigString(viper.GetString("thumbor.base_url"))
	secret := config.GetConfigString(viper.GetString("thumbor.secret"))
	keyForSign := []byte(secret)
	h := hmac.New(sha1.New, keyForSign)
	h.Write([]byte(input))
	replacer := strings.NewReplacer("/", "_", "+", "-")
	signature := replacer.Replace(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	result := fmt.Sprintf("%s/%s/%s", baseUrl, signature, input)
	return result
}

func GetListSize() []map[string]string {
	size_large_screen := config.GetConfigString(viper.GetString("thumbor.size_large_screen"))
	size_small_screen := config.GetConfigString(viper.GetString("thumbor.size_small_screen"))

	var listSize = []map[string]string{
		{"size_small_screen": size_large_screen},
		{"size_large_screen": size_small_screen},
	}

	return listSize
}

func GetNewThumborImages(mediaPath string) datatypes.JSON {
	newMediaPath, _ := url.Parse(mediaPath)
	type mapString map[string]string
	var arrProceedMapString []mapString
	formatImage := config.GetConfigString(viper.GetString("thumbor.format_image"))
	mediaPathThumbor := fmt.Sprintf("%s/%s", formatImage, newMediaPath)
	listSize := GetListSize()
	for _, valSize := range listSize {
		var tmpMapString = make(map[string]string, len(valSize))
		for key, value := range valSize {
			tmpValue := GetThumborUrl(fmt.Sprintf("%s/%s", value, mediaPathThumbor))
			tmpMapString[key] = tmpValue
		}
		arrProceedMapString = append(arrProceedMapString, tmpMapString)
	}
	newMediaImages, _ := json.Marshal(arrProceedMapString)
	return newMediaImages
}
