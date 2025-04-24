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
	switch {
	case errors.Is(err, ErrNotFound):
		http.Error(w, ErrorResponse(err.Error()), http.StatusNotFound)
	case errors.Is(err, ErrBadRequest):
		http.Error(w, ErrorResponse(err.Error()), http.StatusBadRequest)
	case err == nil:
		return
	default:
		http.Error(w, ErrorResponse("Unknown error"), http.StatusInternalServerError)
		h.logger.Errorf("an unknown error occurred in the application: %s", err.Error())
	}
}

func NewHTTPErrorHandler(logger *logrus.Logger) *HTTPErrorHandler {
	return &HTTPErrorHandler{logger: logger}
}
