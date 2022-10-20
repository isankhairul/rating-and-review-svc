package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-klikdokter/helper/message"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/spf13/viper"
	"moul.io/http2curl"
)

var (
	infoLog     = "[API-Payment] Request - HttpRequest.PerformRequest.info"
	responseLog = "[API-Payment] Response - Response.Http.Client"
	errorLog    = "HttpRequest.PerformRequest.error"
)

type OrderExistResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func CheckOrderIdExist(orderId string) (message.Message, error) {
	parameters := make(map[string]interface{})
	parameters["order_id"] = orderId
	jsonData, _ := json.Marshal(parameters)

	url := viper.GetString("payment-service.check-order-id")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		LoggerHttpClient(errorLog, fmt.Sprintf("%v", err))
		return message.ErrFailedRequestToPayment, err
	}

	// Log CURL
	command, _ := http2curl.GetCurlCommand(req)
	LoggerHttpClient(infoLog, fmt.Sprintf("%v", command))

	res, err := client.Do(req)
	if err != nil {
		LoggerHttpClient(errorLog, fmt.Sprintf("%v", err))
		return message.ErrFailedRequestToPayment, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		LoggerHttpClient(errorLog, fmt.Sprintf("%v", err))
		return message.ErrFailedRequestToPayment, err
	}

	// Log Response
	LoggerHttpClient(responseLog, string(body))

	// Condition if response status not 200
	if res.StatusCode != 200 {
		msg := "Error Payment Service Response Status Code is " + strconv.Itoa(res.StatusCode)
		return message.Message{Code: 412002, Message: msg}, nil
	}

	// Check response message that order id is exist
	orderResult := OrderExistResponse{}
	if err := json.Unmarshal([]byte(body), &orderResult); err != nil {
		LoggerHttpClient(errorLog, fmt.Sprintf("%v", err))
		return message.ErrFailedRequestToPayment, err
	}

	if orderResult.Message == "Order id is exist" {
		return message.SuccessMsg, nil
	} else {
		return message.Message{Code: 412002, Message: "Order id is not exist"}, nil
	}
}

func UpdateFlagPayment(orderId string) (message.Message, error) {
	parameters := make(map[string]interface{})
	parameters["order_id"] = orderId
	parameters["is_review"] = true
	jsonData, _ := json.Marshal(parameters)

	url := viper.GetString("payment-service.update-flag")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		LoggerHttpClient(errorLog, fmt.Sprintf("%v", err))
		return message.ErrFailedRequestToPayment, err
	}

	// Log CURL
	command, _ := http2curl.GetCurlCommand(req)
	LoggerHttpClient(infoLog, fmt.Sprintf("%v", command))

	res, err := client.Do(req)
	if err != nil {
		LoggerHttpClient(errorLog, fmt.Sprintf("%v", err))
		return message.ErrFailedRequestToPayment, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		LoggerHttpClient(errorLog, fmt.Sprintf("%v", err))
		return message.ErrFailedRequestToPayment, err
	}

	// Log Response
	LoggerHttpClient(responseLog, string(body))

	if res.StatusCode == 200 {
		return message.SuccessMsg, nil
	} else {
		msg := "Error Payment Service Response Status Code is " + strconv.Itoa(res.StatusCode)
		return message.Message{Code: 412002, Message: msg}, nil
	}
}
