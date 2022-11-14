package errors

type ServiceErrorCode string

const (
	ExternalApiCode  = ServiceErrorCode("err_external_api_error")
	BadGenderCode    = ServiceErrorCode("err_bad_gender")
	BadEmailCode     = ServiceErrorCode("err_bad_email")
	UserNotFoundCode = ServiceErrorCode("err_user_not_found")

	InternalErrorCode = ServiceErrorCode("err_internal")

	ParsingErrorCode = ServiceErrorCode("err_parsing")

	DBErrorSavingCode   = ServiceErrorCode("err_db_save")
	DBErrorDeletingCode = ServiceErrorCode("err_db_delete")
	DBErrorGettingCode  = ServiceErrorCode("err_db_get")
	DBErrorNotFoundCode = ServiceErrorCode("err_db_not_found")
)

type ServiceError struct {
	Code        ServiceErrorCode
	Description string
}

func New(code ServiceErrorCode, description string) ServiceError {
	return ServiceError{Code: code, Description: description}
}

func NewFromError(code ServiceErrorCode, err error) ServiceError {
	return ServiceError{Code: code, Description: err.Error()}
}

func (s ServiceError) Error() string {
	return string(s.Code) + ": " + s.Description
}
