package workouts

import (
	"context"

	repo "github.com/Davidmuthee12/kicker/internals/adapters/postgres/sqlc"
)

type Service interface {
	ListWorkouts(ctx context.Context) ([]repo.Workout, error)
	AddWorkout(ctx context.Context, arg repo.AddWorkoutParams) (repo.Workout, error)
	// GetWorkout(ctx context.Context, id int) (*repo.Workout, error)
	// AddExercise(ctx context.Context, workoutID int, name string) (*repo.Exercise, error)
}

type svc struct {
	// Repository
	repo repo.Querier
}

func (s *svc) AddWorkout(ctx context.Context, arg repo.AddWorkoutParams) (repo.Workout, error) {
	return s.repo.AddWorkout(ctx, arg)
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) ListWorkouts(ctx context.Context) ([]repo.Workout, error) {
	return s.repo.ListWorkouts(ctx)
}