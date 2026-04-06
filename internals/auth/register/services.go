package register

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, input RegisterInput) (User, error)
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	ErrInvalidRegisterInput = errors.New("invalid register input")
	ErrEmailAlreadyExists   = errors.New("email already exists")
)

type svc struct {
	db *pgx.Conn
}

func NewService(db *pgx.Conn) Service {
	return &svc{db: db}
}

func (s *svc) Register(ctx context.Context, input RegisterInput) (User, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.TrimSpace(strings.ToLower(input.Email))

	if name == "" || email == "" || len(input.Password) < 8 {
		return User{}, ErrInvalidRegisterInput
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	createdUser, err := s.createUser(ctx, name, email, string(hashedPassword))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, ErrEmailAlreadyExists
		}

		return User{}, err
	}

	return User{
		ID:    createdUser.ID.String(),
		Name:  createdUser.Name,
		Email: createdUser.Email,
	}, nil
}

type dbUser struct {
	ID    pgtype.UUID
	Name  string
	Email string
}

func (s *svc) createUser(ctx context.Context, name, email, passwordHash string) (dbUser, error) {
	const query = `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, name, email
	`

	var user dbUser
	err := s.db.QueryRow(ctx, query, name, email, passwordHash).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	return user, err
}
