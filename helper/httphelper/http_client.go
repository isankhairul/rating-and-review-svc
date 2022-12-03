package httphelper

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gojek/heimdall"
	"github.com/gojek/heimdall/httpclient"
)

var (
	initialTimeout        = 2 * time.Millisecond // Inital timeout
	maxTimeout            = 9 * time.Millisecond // Max time out
	exponentFactor        = 2                    // Multiplier
	maximumJitterInterval = 2 * time.Millisecond // Max jitter interval. It must be more than 1*time.Millisecond
	timeout               = 15 * time.Second
	retryCounter          = 3
)

//NewClient wrapper
func NewClient() *httpclient.Client {
	// First set a backoff mechanism. Exponential Backoff increases the backoff at a exponential rate
	backoff := heimdall.NewExponentialBackoff(initialTimeout, maxTimeout, float64(exponentFactor), maximumJitterInterval)

	// Create a new retry mechanism with the backoff
	retrier := heimdall.NewRetrier(backoff)

	// Create a new client, sets the retry mechanism, and the number of times you would like to retry
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(retryCounter),
	)

	return client
}

//NewRequest is a wrapper of each request
func NewRequest(method string, url string, body io.ReadSeeker) (*http.Request, error) {
	var reqBody io.ReadCloser

	if body != nil {
		reqBody = ioutil.NopCloser(body)
	}

	httpReq, err := http.NewRequest(method, url, reqBody)

	if err != nil {
		return nil, err
	}

	return httpReq, nil
}
