package httphelper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"moul.io/http2curl"
)

type httpRequest struct {
	Logger log.Logger
}


type HttpRequest interface {
	PerformRequestWithLog(logger log.Logger, method, url string, body []byte, queryParams, headers map[string]string) (int, []byte, error)
}

func PerformRequestMultipartWithLog(logger log.Logger, method, url string, body []byte, queryParams, headers map[string]string, writer *multipart.Writer) (int, []byte, error) {
	fmt.Println("tester =" , url)
	client := NewClient()
	// method, url, body
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	//Headers
	for key, val := range headers {
		req.Header.Add(key, val)
	}

	//Query Params
	q := req.URL.Query()
	for key, val := range queryParams {
		q.Add(key, val)
	}

	// Log Request
	reqId := (strings.Split(uuid.NewString(), "-"))[0]

	http2curl.GetCurlCommand(req)
	// level.Info(logger).Log(fmt.Sprint(reqId, "cURL"), command)

	// Start Request
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	defer resp.Body.Close()

	// Check Response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		level.Info(logger).Log("error", err)
		return 0, nil, err
	}

	//Log Response
	level.Info(logger).Log(fmt.Sprint(reqId, "Resp", resp.StatusCode))
	return resp.StatusCode, data, nil
}