-- name: ListWorkouts :many
SELECT * FROM workouts;

-- name: AddWorkout :one
INSERT INTO workouts (user_id, title, date, notes)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetWorkoutByID :one
SELECT * FROM workouts WHERE id = $1;

-- name: DeleteWorkout :exec
DELETE FROM workouts WHERE id = $1;

-- name: UpdateWorkoutByID :one
UPDATE workouts
SET
	title = $2,
	date = $3,
	notes = $4
WHERE id = $1
RETURNING *;
