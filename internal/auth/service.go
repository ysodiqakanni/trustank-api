package auth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ysodiqakanni/trustank-api/internal/errors"
	"github.com/ysodiqakanni/trustank-api/internal/user"
	"github.com/ysodiqakanni/trustank-api/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Service encapsulates the authentication logic.
type Service interface {
	// authenticate authenticates a user using username and password.
	// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
	Login(ctx context.Context, username string, password string) (string, error)
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() primitive.ObjectID
	// GetName returns the user name.
	GetName() string

	GetRole() []string
}

type service struct {
	signingKey      string
	tokenExpiration int
	logger          log.Logger
	userRepo        user.Repository
}

type LoginRequest struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the CreateAlbumRequest fields.
func (m LoginRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Username, validation.Required, validation.Length(2, 128)),
		validation.Field(&m.Password, validation.Required, validation.Length(5, 128)),
	)
}

// NewService creates a new authentication service.
func NewService(signingKey string, tokenExpiration int, logger log.Logger, userRepo user.Repository) Service {
	return service{signingKey, tokenExpiration, logger, userRepo}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, username, password string) (string, error) {
	if identity := s.authenticate(ctx, username, password); identity != nil {
		return s.generateJWT(identity)
	}
	return "", errors.Unauthorized("")
}

// authenticate authenticates a user using username and password.
// If username and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, username, password string) Identity {
	logger := s.logger.With(ctx, "user", username)

	// first get user by email
	usr, err := s.userRepo.GetByEmail(ctx, username)
	if err != nil {
		logger.Infof("authentication failed")
		return nil
	}

	logger.Infof("user found by email")
	err = bcrypt.CompareHashAndPassword(usr.HashedPassword, []byte(password))
	if err != nil {
		logger.Errorf("authentication failed due to password", err)
		return nil
		// Todo: check what kind of error occurred
		//if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		//	return 0, ErrInvalidCredentials
		//} else {
		//	return 0, err
		//}
	}
	logger.Infof("authentication successful")
	return usr
}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(identity Identity) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   identity.GetID(),
		"name": identity.GetName(),
		"role": identity.GetRole(),
		"exp":  time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(s.signingKey))
}
