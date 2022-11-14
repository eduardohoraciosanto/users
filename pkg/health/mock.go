package health

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) HealthCheck(ctx context.Context) (service, db bool, err error) {
	args := s.Called(ctx)
	return args.Bool(0), args.Bool(1), args.Error(2)
}
