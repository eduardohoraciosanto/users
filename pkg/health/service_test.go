package health

import (
	"context"
	"testing"

	"github.com/eduardohoraciosanto/users/internal/logger"
)

func TestHealthCheck(t *testing.T) {
	service := NewService(logger.NewLogger("health svc unit test", "testing", false))
	s, err := service.HealthCheck(context.TODO())

	if s != true || err != nil {
		t.Errorf("Unexpected values from method: service %t, error %s", s, err)
	}
}
