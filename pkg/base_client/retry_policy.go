package base_client

import (
	"math"
	"net/http"
	"time"
)

type CheckForRetry func(resp *http.Response, err error) (bool, error)

func DefaultRetryPolicy(resp *http.Response, err error) (bool, error) {
	if err != nil {
		return true, err
	}

	if resp.StatusCode == 0 || resp.StatusCode >= 500 {
		return true, nil
	}
	return false, nil
}

type Backoff func(min, max time.Duration, attempts int, resp *http.Response) time.Duration

func DefaultBackoff(min, max time.Duration, attempts int, resp *http.Response) time.Duration {
	multiplication := math.Pow(2, float64(attempts)) * float64(min)
	sleep := time.Duration(multiplication)
	if float64(sleep) != multiplication || sleep > max {
		sleep = max
	}
	return sleep
}
