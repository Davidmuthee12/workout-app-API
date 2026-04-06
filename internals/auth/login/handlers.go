package login

import (
	"errors"
	"log"
	"net/http"

	"github.com/Davidmuthee12/kicker/internals/json"
)

type handler struct {
	service Service
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := json.Read(r.Body, &req); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.Login(r.Context(), LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, ErrInvalidLoginInput) {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		if errors.Is(err, ErrInvalidCredentials) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		log.Println(err)
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, result)
}
