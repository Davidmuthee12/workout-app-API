package workouts

import (
	"log"
	"net/http"

	repo "github.com/Davidmuthee12/kicker/internals/adapters/postgres/sqlc"
	"github.com/Davidmuthee12/kicker/internals/json"
	"github.com/jackc/pgx/v5/pgtype"
)

type handler struct {
	service Service
}

type createWorkoutRequest struct {
	Title string `json:"title"`
	Date  string `json:"date"`
	Notes string `json:"notes"`
}

type addExerciseRequest struct {
	Name string `json:"name"`
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h handler) ListWorkouts(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.service.ListWorkouts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to list workouts", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, workouts)
}

func (h handler) AddWorkout(w http.ResponseWriter, r *http.Request) {
	var req createWorkoutRequest

	// DECODE REQUEST BODY
	if err := json.Read(r.Body, &req); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// BASIC VALIDATION
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	userIDValue := r.Context().Value("userID")
	userIDString, ok := userIDValue.(string)
	if !ok || userIDString == "" {
		http.Error(w, "Missing user id in token", http.StatusUnauthorized)
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(userIDString); err != nil {
		http.Error(w, "Invalid user id in token", http.StatusBadRequest)
		return
	}

	params := repo.AddWorkoutParams{
		UserID: userID,
		Title:  req.Title,
		Notes:  pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	}

	if req.Date != "" {
		if err := params.Date.Scan(req.Date); err != nil {
			http.Error(w, "Date must be in YYYY-MM-DD format", http.StatusBadRequest)
			return
		}
	}

	workout, err := h.service.AddWorkout(r.Context(), params)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, workout)
}