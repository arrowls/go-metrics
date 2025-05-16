package apperrors

import (
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

type HTTPErrorHandler struct {
	logger *logrus.Logger
}

func (h *HTTPErrorHandler) Handle(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var status int

	switch {
	case err == nil:
		return
	case errors.Is(err, ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, ErrBadRequest):
		status = http.StatusBadRequest
	default:
		status = http.StatusInternalServerError
		h.logger.Errorf("an unknown error occurred in the application: %s", err.Error())
	}

	w.WriteHeader(status)

	_, errWrite := w.Write(ErrorResponse(err.Error()))
	if errWrite != nil {
		h.logger.Error(errWrite)
		return
	}
}

func NewHTTPErrorHandler(logger *logrus.Logger) *HTTPErrorHandler {
	return &HTTPErrorHandler{logger: logger}
}
