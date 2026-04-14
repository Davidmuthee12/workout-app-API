package programs

import (
	"crypto/rand"
	"log"
	"net/http"

	repo "github.com/Davidmuthee12/kicker/internals/adapters/postgres/sqlc"
	"github.com/Davidmuthee12/kicker/internals/json"
	"github.com/jackc/pgx/v5/pgtype"
)

type handler struct {
	service Service
}

type createProgramRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type addProgramRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h handler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	programs, err := h.service.ListPrograms(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to list programs", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, programs)
}

func (h handler) AddProgram(w http.ResponseWriter, r *http.Request) {
	var req addProgramRequest
	// DECODE REQUEST BODY
	if err := json.Read(r.Body, &req); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// BASIC VALIDATION
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	var programID pgtype.UUID
	if _, err := rand.Read(programID.Bytes[:]); err != nil {
		log.Println(err)
		http.Error(w, "Failed to generate program id", http.StatusInternalServerError)
		return
	}
	programID.Valid = true

	program, err := h.service.AddProgram(r.Context(), repo.AddProgramParams{
		ProgramID:   programID,
		Title:       req.Name,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to add program", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, program)
}