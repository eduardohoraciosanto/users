package response

import "fmt"

const (
	ErrCodeInternalServerError        = "err_internal"
	ErrDescriptionInternalServerError = "Internal Server Error"

	ErrCodeBadRequest            = "err_bad_request"
	ErrDescriptionBadRequestURL  = "The URL In Request contains errors"
	ErrDescriptionBadRequestBody = "The provided body contains errors"
)

var (
	StandardInternalServerError = Error{ErrCodeInternalServerError, ErrDescriptionInternalServerError}
	StandardBadBodyRequest      = Error{ErrCodeBadRequest, ErrDescriptionBadRequestBody}
)

type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s:%s", e.Code, e.Description)
}
