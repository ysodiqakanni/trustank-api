package user

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/ysodiqakanni/trustank-api/internal/entity"
	"github.com/ysodiqakanni/trustank-api/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Service encapsulates use case logic for businessCategories.
type Service interface {
	Get(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, req CreateUserRequest) (*User, error)
}

// User represents the data about a User.
type User struct {
	entity.User
}

type CreateUserRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

// Validate validates the CreateAlbumRequest fields.
func (m CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Email, validation.Required, is.Email, validation.Length(6, 200)),
		validation.Field(&m.Password, validation.Required, validation.Length(6, 100)),

		//validation.Field(&a.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new category service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the album with the specified the album ID.
func (s service) Get(ctx context.Context, id primitive.ObjectID) (*User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &User{user}, nil
}

func (s service) GetByEmail(ctx context.Context, email string) (User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}
	return User{user}, nil
}
func (s service) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	existing, getErr := s.GetByEmail(ctx, req.Name)
	//emptyObj := User{}

	if getErr != nil && existing.ID != primitive.NewObjectID() {
		return nil, errors.New("A user with this name already exists")
	}

	//password :=
	// Todo: generate user password
	now := time.Now()
	id, err := s.repo.Create(ctx, entity.User{
		Name:      req.Name,
		Email:     req.Email,
		Role:      req.Roles,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, *id)
}
