package error

import (
	"errors"
	"fmt"
)

type ResponseError struct {
	Message    string
	StatusCode int
	Url        string
	Err        error
}

func NewRequestError(url string, code int, err error, msg string) error {
	return ResponseError{
		Message:    msg,
		StatusCode: code,
		Url:        url,
		Err:        err,
	}
}

// Unwrap returns inner error
func (err ResponseError) Unwrap() error {
	return err.Err
}

func (err ResponseError) Error() string {
	if err.Err != nil {
		return fmt.Sprintf("%v: %v", err.Message, err.Err)
	}
	return err.Message
}

// Dig returns the most nested CustomErrorWrapper
func (err ResponseError) Dig() ResponseError {
	var ew ResponseError
	if errors.As(err.Err, &ew) {
		// Recursively digs until wrapper error is not in which case it will stop
		return ew.Dig()
	}
	return err
}
