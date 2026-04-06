package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	repo "github.com/Davidmuthee12/kicker/internals/adapters/postgres/sqlc"
	"github.com/Davidmuthee12/kicker/internals/auth/login"
	"github.com/Davidmuthee12/kicker/internals/auth/register"
	"github.com/Davidmuthee12/kicker/internals/workouts"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config config
	db     *pgx.Conn
}

type config struct {
	addr      string
	db        dbConfig
	jwtSecret string
}

type dbConfig struct {
	dsn string
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Mount

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	queries := repo.New(app.db)

	// A good base middleware stack
	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // important for rate limiting
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	// PUBLIC ROUTES
	registerService := register.NewService(app.db)
	registerHandler := register.NewHandler(registerService)
	loginService := login.NewService(app.db, app.config.jwtSecret)
	loginHandler := login.NewHandler(loginService)
	r.Post("/auth/register", registerHandler.Register)
	r.Post("/auth/login", loginHandler.Login)

	// PROTECTED ROUTES
	workoutServices := workouts.NewService(queries)
	workoutHandler := workouts.NewHandler(workoutServices)
	r.Route("/api", func(r chi.Router) {
		r.Use(app.authMiddleware)

		r.Get("/workouts", workoutHandler.ListWorkouts)
		r.Post("/workouts", workoutHandler.AddWorkout)
		r.Get("/workouts/{id}", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "not implemented", http.StatusNotImplemented)
		})
		r.Post("/workouts/{id}/exercises", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "not implemented", http.StatusNotImplemented)
		})
	})

	return r
}

// Run
func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server is running at address %s", app.config.addr)
	return srv.ListenAndServe()
}

// AUTH MIDDLEWARE
func (app *application) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(app.config.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
