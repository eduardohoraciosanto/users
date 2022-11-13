package response_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	serviceErrors "github.com/eduardohoraciosanto/users/internal/errors"
	"github.com/eduardohoraciosanto/users/internal/response"
	"github.com/stretchr/testify/assert"
)

func TestRespondWithData(t *testing.T) {
	rec := httptest.NewRecorder()

	err := response.RespondWithData(rec, http.StatusTeapot, "TestData")
	assert.Nil(t, err)

	res := rec.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusTeapot, res.StatusCode)
	assert.Contains(t, string(data), "TestData")

}

func TestRespondWithError_InternalServer(t *testing.T) {
	rec := httptest.NewRecorder()

	err := response.RespondWithError(rec, fmt.Errorf("some Error"))
	assert.Nil(t, err)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestRespondWithError_IsServiceError_InternalError(t *testing.T) {
	rec := httptest.NewRecorder()

	serviceError := serviceErrors.ServiceError{
		Code: "someCode",
	}

	err := response.RespondWithError(rec, serviceError)
	assert.Nil(t, err)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestRespondWithError_IsError_InternalError(t *testing.T) {
	rec := httptest.NewRecorder()

	resError := response.Error{
		Code:        "someCode",
		Description: "someDescription",
	}

	err := response.RespondWithError(rec, resError)
	assert.Nil(t, err)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestRespondWithError_IsError_BadRequest(t *testing.T) {
	rec := httptest.NewRecorder()

	resError := response.Error{
		Code:        response.ErrCodeBadRequest,
		Description: response.ErrDescriptionBadRequestBody,
	}

	err := response.RespondWithError(rec, resError)
	assert.Nil(t, err)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}
