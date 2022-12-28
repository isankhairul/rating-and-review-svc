package thumbor

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"go-klikdokter/helper/config"
	"net/url"
	"strings"

	"github.com/spf13/viper"
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

func GetNewThumborImages(mediaPath string) string {
	newMediaPath, _ := url.Parse(mediaPath)
	formatImage := config.GetConfigString(viper.GetString("thumbor.format_image"))
	sizeLargeScreen := config.GetConfigString(viper.GetString("thumbor.size_large_screen"))
	mediaPathThumbor := fmt.Sprintf("%s/%s/%s", sizeLargeScreen, formatImage, newMediaPath)
	result := GetThumborUrl(mediaPathThumbor)

	return result
}

func GetNewThumborImagesOriginal(mediaPath string) string {
	newMediaPath, _ := url.Parse(mediaPath)
	formatImage := config.GetConfigString(viper.GetString("thumbor.format_image"))
	sizeOriginal := config.GetConfigString(viper.GetString("thumbor.size_ar_original"))
	mediaPathThumbor := fmt.Sprintf("%s/%s/%s", sizeOriginal, formatImage, newMediaPath)
	result := GetThumborUrl(mediaPathThumbor)

	return result
}
