package health

import (
	"context"
	"testing"

	"github.com/eduardohoraciosanto/users/internal/db"
	"github.com/eduardohoraciosanto/users/internal/logger"
	"github.com/stretchr/testify/mock"
)

func TestHealthCheck(t *testing.T) {
	dbMock := &db.DBMock{}

	dbMock.On("Alive", mock.Anything).Return(true)

	service := NewService(logger.NewLogger("health svc unit test", "testing", false), dbMock)
	s, d, err := service.HealthCheck(context.TODO())

	if s != true || d != true || err != nil {
		t.Errorf("Unexpected values from method: service %t, error %s", s, err)
	}
}
