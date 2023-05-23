package base_client

import (
	"math"
	"math/rand"
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

type Backoff func(min, max int, attempts int) time.Duration

func DefaultBackoff(min, max int, attempts int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	m := int(math.Pow(2, float64(attempts)) - 1)

	sleep := r.Intn(m-min) + min
	if sleep > max {
		sleep = max
	}

	return time.Duration(sleep) * time.Second
}
