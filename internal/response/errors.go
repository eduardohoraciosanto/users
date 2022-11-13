package response

import "fmt"

const (
	ErrCodeInternalServerError        = "err_internal"
	ErrDescriptionInternalServerError = "Internal Server Error"

	ErrCodeBadRequest            = "err_bad_request"
	ErrDescriptionBadRequestURL  = "The URL In Request contains errors"
	ErrDescriptionBadRequestBody = "The provided body contains errors"

	ErrDescriptionCartNotFound = "The Cart ID was not found"

	ErrDescriptionItemAlreadyInCart = "The item already exists in the cart"
	ErrDescriptionItemNotFound      = "The item does not exists in the cart"

	ErrDescriptionItemNotFoundProvider = "The item was not found on the provider"
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
