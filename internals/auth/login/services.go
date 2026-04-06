package login

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(ctx context.Context, input LoginInput) (LoginResult, error)
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	ErrInvalidLoginInput  = errors.New("invalid login input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrMissingJWTSecret   = errors.New("missing jwt secret")
)

type svc struct {
	db        *pgx.Conn
	jwtSecret string
	tokenTTL  time.Duration
}

func NewService(db *pgx.Conn, jwtSecret string) Service {
	return &svc{
		db:        db,
		jwtSecret: jwtSecret,
		tokenTTL:  24 * time.Hour,
	}
}

func (s *svc) Login(ctx context.Context, input LoginInput) (LoginResult, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return LoginResult{}, ErrInvalidLoginInput
	}

	if s.jwtSecret == "" {
		return LoginResult{}, ErrMissingJWTSecret
	}

	user, err := s.getUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return LoginResult{}, ErrInvalidCredentials
		}

		return LoginResult{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	now := time.Now()
	userID := user.ID.String()
	claims := jwt.MapClaims{
		"user_id": userID,
		"sub":     userID,
		"iat":     now.Unix(),
		"exp":     now.Add(s.tokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		Token: signedToken,
		User: User{
			ID:    userID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

type dbUser struct {
	ID           pgtype.UUID
	Name         string
	Email        string
	PasswordHash string
}

func (s *svc) getUserByEmail(ctx context.Context, email string) (dbUser, error) {
	const query = `
		SELECT id, name, email, password_hash
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	var user dbUser
	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)

	return user, err
}
