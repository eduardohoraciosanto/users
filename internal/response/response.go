package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/eduardohoraciosanto/users/internal/config"
	serviceErrors "github.com/eduardohoraciosanto/users/internal/errors"
)

type Meta struct {
	Version string `json:"version"`
}

type BaseResponse struct {
	Meta  Meta        `json:"meta"`
	Data  interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

func newBaseResponseWithData(data interface{}) BaseResponse {
	return BaseResponse{
		Meta: Meta{
			Version: config.GetVersion(),
		},
		Data: data,
	}
}

func newBaseResponseWithError(err interface{}) BaseResponse {
	return BaseResponse{
		Meta: Meta{
			Version: config.GetVersion(),
		},
		Error: err,
	}
}

func RespondWithData(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(newBaseResponseWithData(data))
}

func RespondWithError(w http.ResponseWriter, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeFromError(err))
	return json.NewEncoder(w).Encode(newBaseResponseWithError(viewModelFromError(err)))
}

func statusCodeFromError(err error) int {
	mErr := &serviceErrors.ServiceError{}
	if errors.As(err, mErr) {
		switch mErr.Code {
		default:
			return http.StatusInternalServerError
		}
	}
	vErr := &Error{}
	if errors.As(err, vErr) {
		switch vErr.Code {
		case ErrCodeBadRequest:
			return http.StatusBadRequest
		default:
			return http.StatusInternalServerError
		}
	}

	return http.StatusInternalServerError
}

func descriptionFromError(mErr *serviceErrors.ServiceError) string {

	switch mErr.Code {
	}
	return ErrDescriptionInternalServerError
}

func viewModelFromError(err error) Error {
	sErr := &serviceErrors.ServiceError{}
	if errors.As(err, sErr) {
		return Error{
			Code:        err.Error(),
			Description: descriptionFromError(sErr),
		}
	}
	vErr := Error{}
	if errors.As(err, &vErr) {
		return vErr
	}
	return StandardInternalServerError
}
