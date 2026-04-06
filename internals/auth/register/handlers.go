package register

import (
	"errors"
	"log"
	"net/http"

	"github.com/Davidmuthee12/kicker/internals/json"
)

type handler struct {
	service Service
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	if err := json.Read(r.Body, &req); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, ErrInvalidRegisterInput) {
			http.Error(w, "Name, email and password (min 8 chars) are required", http.StatusBadRequest)
			return
		}

		if errors.Is(err, ErrEmailAlreadyExists) {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}

		log.Println(err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, user)
}
