package api

import (
	"net/http"
)

type StatusError struct {
	Code    int
	Message string
}

func (e StatusError) Error() string {
	return e.Message
}

var (
	NotFoundError      = StatusError{Message: "Resource not found", Code: http.StatusNotFound}
	AuthKeyError       = StatusError{Message: "Auth Key is missing", Code: http.StatusBadRequest}
	InvalidApiKeyError = StatusError{Message: "Requires a valid api key", Code: http.StatusBadRequest}
	ParamsMissingError = StatusError{Message: "Required parameters are missing", Code: http.StatusBadRequest}
	KeenAddonError     = StatusError{Message: "Keen Addon Invocation failed", Code: http.StatusBadRequest}
)
