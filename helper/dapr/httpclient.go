package helper_dapr

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"go-klikdokter/helper/config"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpClient interface {
	PublishEvent(topic string, data string) (map[string]interface{}, error)
}

type httpClient struct {
	Host         string
	Port         string
	Version      string
	PubSubName   string
	Topic        string
	Data         []byte
	IsRawPayload bool
}

func NewDaprHttpClient() HttpClient {
	var host string = config.GetConfigString(viper.GetString("dapr.host"))
	var port string = config.GetConfigString(viper.GetString("dapr.port"))
	var version string = config.GetConfigString(viper.GetString("dapr.version"))
	var pubsubName string = config.GetConfigString(viper.GetString("dapr.pubsub-name"))

	return &httpClient{
		Host:         host,
		Port:         port,
		Version:      version,
		PubSubName:   pubsubName,
		IsRawPayload: true,
	}
}

func (c httpClient) PublishEvent(topic string, data string) (map[string]interface{}, error) {
	client := http.Client{}
	url := fmt.Sprintf("%s:%s/%s/publish/%s/%s", c.Host, c.Port, c.Version, c.PubSubName, topic)
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	// Publish an event using Dapr pub/sub
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var dataResponse = make(map[string]interface{})
	dataResponse = map[string]interface{}{
		"URL":            url,
		"Request":        "",
		"ResponseHeader": "",
		"Response":       "",
	}

	dataResponse["Request"] = string(data)
	dataRsh, _ := json.Marshal(response.Header)
	dataRs, _ := ioutil.ReadAll(response.Body)
	dataResponse["ResponseHeader"] = string(dataRsh)
	dataResponse["Response"] = string(dataRs)

	return dataResponse, nil
}
