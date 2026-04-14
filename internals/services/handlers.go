package services

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

type createServiceRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type addServiceRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h handler) ListServices(w http.ResponseWriter, r *http.Request) {
	services, err := h.service.ListServices(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to list services", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, services)
}

func (h handler) AddService(w http.ResponseWriter, r *http.Request) {
	var req addServiceRequest
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

	var serviceID pgtype.UUID
	if _, err := rand.Read(serviceID.Bytes[:]); err != nil {
		log.Println(err)
		http.Error(w, "Failed to generate service id", http.StatusInternalServerError)
		return
	}
	serviceID.Valid = true

	service, err := h.service.AddService(r.Context(), repo.AddServiceParams{
		ServiceID:   serviceID,
		Title:       req.Name,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
	})

	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to add service", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, service)
}