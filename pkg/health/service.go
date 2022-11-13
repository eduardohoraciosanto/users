package health

import (
	"context"

	"github.com/eduardohoraciosanto/users/internal/logger"
)

//Service is the interface for the health
type Service interface {
	HealthCheck(ctx context.Context) (service bool, err error)
}

type svc struct {
	log logger.Logger
}

//NewService gives a new Service
func NewService(log logger.Logger) Service {
	return &svc{
		log: log,
	}
}

//HealthCheck returns the status of the API and it's components
func (s *svc) HealthCheck(ctx context.Context) (service bool, err error) {
	s.log.Info(ctx, "Performing Health Check")

	return true, nil
}
