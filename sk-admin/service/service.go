package service

type Service interface {
	HealthCheck() bool
}

type SkAdminService struct {
}

func (s SkAdminService) HealthCheck() bool {
	return true
}

type ServiceMiddleware func(service Service) Service
