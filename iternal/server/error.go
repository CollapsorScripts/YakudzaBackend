package server

import (
	"Yakudza/pkg/logger"
	"Yakudza/pkg/utilities"
	"net/http"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SetHTTPError(w http.ResponseWriter, errStr string, code int) {
	w.WriteHeader(code)

	h := HTTPError{
		Code:    code,
		Message: errStr,
	}

	str := utilities.ToJSON(h)
	_, err := w.Write([]byte(str))
	if err != nil {
		logger.Error("%v", err)
	}
}
