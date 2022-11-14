package health_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eduardohoraciosanto/users/pkg/health"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHealth_OK(t *testing.T) {

	sMock := &health.ServiceMock{}
	sMock.On("HealthCheck", mock.Anything).Return(true, true, nil)

	h := health.Handler{
		Service: sMock,
	}

	req, err := http.NewRequest("GET", "/health", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.Health(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestHealth_ServiceError(t *testing.T) {
	sMock := &health.ServiceMock{}
	sMock.On("HealthCheck", mock.Anything).Return(true, true, errors.New("mock was asked to fail"))

	h := health.Handler{
		Service: sMock,
	}

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.Health(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
