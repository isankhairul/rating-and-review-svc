package httputil

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var (
	defaultWaitMin    = 1 * time.Second
	defaultWaitMax    = 20 * time.Second
	defaultAttemptMax = 2
)

//Request model
type Request struct {
	//body is use to rewind the request data between retires
	body io.ReadSeeker

	*http.Request
}

//Client model
type Client struct {
	HTTPClient      *http.Client
	RetryWaitMin    time.Duration
	RetryWaitMax    time.Duration
	RetryAttemptMax int
	Retry           Retry
	BackOff         BackOff
}

//NewClient wrapper
func NewClient() *Client {
	defaultTimeout := 30 * time.Second
	defaultTimeoutHandshake := 15 * time.Second
	var netTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: defaultTimeoutHandshake,
		}).DialContext,
		TLSHandshakeTimeout: defaultTimeoutHandshake,
	}
	return &Client{
		HTTPClient: &http.Client{
			Timeout:   defaultTimeout,
			Transport: netTransport,
		},
		RetryWaitMin:    defaultWaitMin,
		RetryWaitMax:    defaultWaitMax,
		RetryAttemptMax: defaultAttemptMax,
		Retry:           RetryPolicy,
		BackOff:         BackOffPolicy,
	}
}

//NewRequest is a wrapper of each request
func NewRequest(method string, url string, body io.ReadSeeker) (*Request, error) {
	var reqBody io.ReadCloser

	if body != nil {
		reqBody = ioutil.NopCloser(body)
	}

	httpReq, err := http.NewRequest(method, url, reqBody)

	if err != nil {
		return nil, err
	}

	return &Request{body, httpReq}, nil
}

//Do is a wrapper for calling http method with retries
func (c *Client) Do(req *Request) (*http.Response, error) {
	retryAttempt := 0
	for {
		if req.body != nil {
			if _, err := req.body.Seek(0, 0); err != nil {
				return nil, err
			}
		}

		//Attempt to request
		resp, err := c.HTTPClient.Do(req.Request)

		//Check if we need retry the request
		ok, checkErr := c.Retry(resp, err)

		if err != nil {
			fmt.Printf("%s %s request failed: %v ", req.Method, req.URL, err)
		}

		if !ok {
			if checkErr != nil {
				err = checkErr
			}

			return resp, err
		}

		//Check if max attempt value is
		if retryAttempt == c.RetryAttemptMax {
			break
		}

		wait := c.BackOff(c.RetryWaitMin, c.RetryWaitMax, retryAttempt, resp)
		retryAttempt++ //add attempt
		fmt.Printf("%s %s: retrying in %v seconds, attempt: %d", req.Method, req.URL, wait, retryAttempt)
		time.Sleep(wait)
	}

	return nil, fmt.Errorf("Request from %s %s is failed after %d attempts", req.Method, req.URL, c.RetryAttemptMax)
}
