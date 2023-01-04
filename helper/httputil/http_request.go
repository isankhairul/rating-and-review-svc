package httputil

import (
	"bytes"
	"encoding/json"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/viper"
	"go-klikdokter/helper/config"
	"gorm.io/gorm"
	"io/ioutil"
	logInternal "log"
	"moul.io/http2curl"
	"net/http"
	"os"
	"time"
)

type HttpRequestStruct struct {
	Request      *http.Request
	Response     *http.Response
	Db           *gorm.DB
	Logger       log.Logger
	ErrorMessage string
}

type HttpRequest interface {
	PerformRequest(
		method string,
		url string,
		body []byte,
		bodyQuery []byte,
		headers map[string]string,
		Logger log.Logger,
	) (int, []byte, map[string]interface{}, error)
}

func NewHttpRequest() HttpRequest {
	return &HttpRequestStruct{}
}

func (hr *HttpRequestStruct) PerformRequest(
	method string,
	url string,
	body []byte,
	bodyQuery []byte,
	headers map[string]string,
	Logger log.Logger,
) (int, []byte, map[string]interface{}, error) {
	//errorLog := "HttpRequest.PerformRequest.error"
	//infoLog := "HttpRequest.PerformRequest.info"
	logUID, _ := gonanoid.New()
	client := NewClient()

	var dataResponse = make(map[string]interface{})
	dataResponse = map[string]interface{}{
		"URL":            url,
		"RequestHeader":  "",
		"Request":        "",
		"ResponseHeader": "",
		"Response":       "",
	}

	var req *Request
	var err error
	if method == "GET" {
		req, err = NewRequest(method, url, bytes.NewReader(body))
		command, _ := http2curl.GetCurlCommand(req.Request)
		if err != nil {
			dataResponse["Response"] = "Error when request get, with message: " + err.Error()
			//LoggerHttpClient(errorLog, fmt.Sprintf("%v", err))
			_ = level.Error(Logger).Log("logUID", logUID, "RQ", string(body), "RS", err.Error(), "Curl", command)
			return 0, nil, dataResponse, err
		}
		req.Header.Add("Content-Type", "application/json")
		req.URL.RawQuery = string(bodyQuery)
	} else {
		req, err = NewRequest(method, url, bytes.NewReader(body))
		command, _ := http2curl.GetCurlCommand(req.Request)
		if err != nil {
			dataResponse["Response"] = "Error when request, with message: " + err.Error()
			//LoggerHttpClient(errorLog, fmt.Sprintf("%v", err))
			_ = level.Error(Logger).Log("logUID", logUID, "RQ", string(body), "RS", err.Error(), "Curl", command)
			return 0, nil, dataResponse, err
		}
		req.Header.Add("Content-Type", "application/json")
	}

	for key, val := range headers {
		req.Header.Add(key, val)
	}
	jsonHeaders, _ := json.Marshal(headers)

	dataResponse["RequestHeader"] = string(jsonHeaders)
	dataResponse["Request"] = string(body)

	//Log CURL
	command, _ := http2curl.GetCurlCommand(req.Request)
	//LoggerHttpClient(infoLog, fmt.Sprintf("%v", command))

	resp, err := client.Do(req)

	if err != nil {
		dataResponse["Response"] = "Error when client do, with message: " + err.Error()
		_ = level.Error(Logger).Log("logUID", logUID, "RQ", string(body), "RS", err.Error(), "Curl", command)
		return 0, nil, dataResponse, err
	}

	dataRsh, err := json.Marshal(resp.Header)
	if err != nil {
		dataResponse["Response"] = "Error when jsonMarshal resp.Header, with message: " + err.Error()
		_ = level.Error(Logger).Log("logUID", logUID, "RQ", string(body), "RS", err.Error(), "Curl", command)
		return 0, nil, dataResponse, err
	}

	dataRs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		dataResponse["Response"] = "Error when ioutil.ReadAll body, with message: " + err.Error()
		_ = level.Error(Logger).Log("logUID", logUID, "RQ", string(body), "RS", err.Error(), "Curl", command)
		//LoggerHttpClient(errorLog, fmt.Sprintf("%v", err))
		return 0, nil, dataResponse, err
	}
	_ = level.Debug(Logger).Log("logUID", logUID, "RQ", string(body), "RS", string(dataRs), "Curl", command)
	dataResponse["ResponseHeader"] = string(dataRsh)
	dataResponse["Response"] = string(dataRs)

	defer resp.Body.Close()

	return resp.StatusCode, dataRs, dataResponse, err
}

type loggingResponseWriter struct {
	status int
	body   interface{}
	http.ResponseWriter
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(body []byte) (int, error) {
	var data interface{}
	_ = json.Unmarshal(body, &data)
	w.body = data
	return w.ResponseWriter.Write(body)
}

func LoggerRequestResponse(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathFolder := config.GetConfigString(viper.GetString("server.output-request-response-logging-path"))

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logInternal.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		// Work / inspect body. You may even modify it!
		var data interface{}
		_ = json.Unmarshal(body, &data)
		reqMap := map[string]interface{}{
			"host":       r.Host,
			"requestUri": r.RequestURI,
			"header":     r.Header,
			"method":     r.Method,
			"body":       data,
			//"summary": a,
		}

		LoggingToFile(pathFolder, "[Request]", reqMap)

		// And now set a new body, which will simulate the same data we read:
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		// Create a response wrapper:
		loggingRW := &loggingResponseWriter{
			ResponseWriter: w,
		}
		// Call next handler, passing the response wrapper:
		h.ServeHTTP(loggingRW, r)
		resMap := map[string]interface{}{
			"status": loggingRW.status,
			"body":   loggingRW.body,
		}

		LoggingToFile(pathFolder, "[Response]", resMap)
	})
}

func LoggerHttpClient(logType string, input interface{}) {
	pathFolder := config.GetConfigString(viper.GetString("server.output-logging-path"))

	LoggingToFile(pathFolder, logType, input)
}

func LoggingToFile(pathFolder string, logType string, input interface{}) {
	pathFileName := pathFolder + time.Now().Format("20060102") + ".log"

	_ = os.MkdirAll(pathFolder, os.ModePerm)

	f, err := os.OpenFile(pathFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logInternal.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	inputSerialize, _ := json.Marshal(input)

	logInternal.SetOutput(f)
	logInternal.Println(logType, string(inputSerialize))
}
