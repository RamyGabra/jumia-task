package custom_errors

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type CustomError struct {
	err        error
	message    string
	StatusCode int
	service    string
	method     string
}

func (e *CustomError) Error() string {
	return e.message
}

func (e *CustomError) Unwrap() error {
	return e.err
}

func (e *CustomError) Log() {

	logrus.WithFields(logrus.Fields{
		"service": e.service,
	}).Errorf("Error %s - %s", e.method, e.Unwrap())

}

func NewInternalServerError(err error, service string, method string) *CustomError {
	return &CustomError{
		err:        err,
		message:    "Internal Server Error: please try again later",
		StatusCode: http.StatusInternalServerError,
		service:    service,
		method:     method,
	}
}

func NewBadRequestError(err error, message string, service string, method string) *CustomError {
	return &CustomError{
		err:        err,
		message:    "Bad request Error: " + message,
		StatusCode: http.StatusBadRequest,
		service:    service,
		method:     method,
	}
}
