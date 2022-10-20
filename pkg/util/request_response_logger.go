package util

import (
	"encoding/json"
	"go-klikdokter/helper/config"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

func LoggerHttpClient(logType string, input interface{}) {
	pathFolder := config.GetConfigString(viper.GetString("server.output-payment-logging-path"))

	LoggingToFile(pathFolder, logType, input)
}

func LoggingToFile(pathFolder string, logType string, input interface{}) {
	pathFileName := pathFolder + time.Now().Format("20060102") + ".log"

	_ = os.MkdirAll(pathFolder, os.ModePerm)

	f, err := os.OpenFile(pathFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	inputSerialize, _ := json.Marshal(input)

	log.SetOutput(f)
	log.Println(logType, string(inputSerialize))
}
