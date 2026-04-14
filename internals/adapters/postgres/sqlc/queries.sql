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

-- name: ListPrograms :many
SELECT * FROM programs;

-- name: GetProgramByID :one
SELECT * FROM programs WHERE id = $1;

-- name: AddProgram :one
INSERT INTO programs (program_id, title, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListServices :many
SELECT * FROM services;

-- name: GetServiceByID :one
SELECT * FROM services WHERE id = $1;

-- name: AddService :one
INSERT INTO services (service_id, title, description)
VALUES ($1, $2, $3)
RETURNING *;

