package errors

const (
	ExternalApiErrorCode = "err_external_api_error"
)

type ServiceError struct {
	Code string
}

func (s ServiceError) Error() string {
	return s.Code
}
