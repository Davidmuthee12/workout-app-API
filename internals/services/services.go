package services

import (
	"context"

	repo "github.com/Davidmuthee12/kicker/internals/adapters/postgres/sqlc"
)

type Service interface {
	ListServices(ctx context.Context) ([]repo.Service, error)
	AddService(ctx context.Context, arg repo.AddServiceParams) (repo.Service, error)
}

type svc struct {
	// Repository
	repo repo.Querier
}

func (s *svc) AddService(ctx context.Context, arg repo.AddServiceParams) (repo.Service, error) {
	return s.repo.AddService(ctx, arg)
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) ListServices(ctx context.Context) ([]repo.Service, error) {
	return s.repo.ListServices(ctx)
}