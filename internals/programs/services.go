package programs

import (
	"context"

	repo "github.com/Davidmuthee12/kicker/internals/adapters/postgres/sqlc"
)

type Service interface {
	ListPrograms(ctx context.Context) ([]repo.Program, error)
	AddProgram(ctx context.Context, arg repo.AddProgramParams) (repo.Program, error)
	// GetProgram(ctx context.Context, id int) (*repo.Program, error)
	// AddExercise(ctx context.Context, programID int, name string) (*repo.Exercise, error)
}

type svc struct {
	// Repository
	repo repo.Querier
}

func (s *svc) AddProgram(ctx context.Context, arg repo.AddProgramParams) (repo.Program, error) {
	return s.repo.AddProgram(ctx, arg)
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) ListPrograms(ctx context.Context) ([]repo.Program, error) {
	return s.repo.ListPrograms(ctx)
}