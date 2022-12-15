package httputil

import (
	"math"
	"net/http"
	"time"
)

//Retry is used for specify policies for handling retries,
//If return false the client stop retrying and return response to the caller,
//If return an error, the error value returned in lieu from the request
type Retry func(resp *http.Response, err error) (bool, error)

//BackOff is used to specify policy how long to wait between retry
type BackOff func(min time.Duration, max time.Duration, retryAttempt int, resp *http.Response) time.Duration

//RetryPolicy provide a callback for Retry client
func RetryPolicy(resp *http.Response, err error) (bool, error) {
	if err != nil {
		return true, err
	}

	if resp.StatusCode == 0 || resp.StatusCode == 500 {
		return true, nil
	}

	return false, nil
}

//BackOffPolicy provice backoff callback for client
func BackOffPolicy(min time.Duration, max time.Duration, retryAttempt int, resp *http.Response) time.Duration {
	//Exponential waiting time between retry
	waitingTime := time.Duration(math.Pow(2, float64(retryAttempt)) * float64(min))

	if waitingTime > max {
		waitingTime = max
	}

	return waitingTime
}
