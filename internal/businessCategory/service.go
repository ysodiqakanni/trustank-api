package businessCategory

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ysodiqakanni/trustank-api/internal/entity"
	"github.com/ysodiqakanni/trustank-api/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Service encapsulates use case logic for businessCategories.
type Service interface {
	Get(ctx context.Context, id primitive.ObjectID) (BusinessCategory, error)
	GetByName(ctx context.Context, name string) (BusinessCategory, error)
	Create(ctx context.Context, req CreateBusinessCategoryRequest) (BusinessCategory, error)
}

// BusinessCategory represents the data about a BusinessCategory.
type BusinessCategory struct {
	entity.BusinessCategory
}

// CreateBusinessCategoryRequest represents an category creation request.
type CreateBusinessCategoryRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreateAlbumRequest fields.
func (m CreateBusinessCategoryRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),

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
func (s service) Get(ctx context.Context, id primitive.ObjectID) (BusinessCategory, error) {
	category, err := s.repo.Get(ctx, id)
	if err != nil {
		return BusinessCategory{}, err
	}
	return BusinessCategory{category}, nil
}

func (s service) GetByName(ctx context.Context, name string) (BusinessCategory, error) {
	category, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return BusinessCategory{}, err
	}
	return BusinessCategory{category}, nil
}
func (s service) Create(ctx context.Context, req CreateBusinessCategoryRequest) (BusinessCategory, error) {
	if err := req.Validate(); err != nil {
		return BusinessCategory{}, err
	}

	existing, _ := s.GetByName(ctx, req.Name)
	emptyObj := BusinessCategory{}
	if existing != emptyObj {
		return BusinessCategory{}, errors.New("A business_ category with this name already exists")
	}

	now := time.Now()
	id, err := s.repo.Create(ctx, entity.BusinessCategory{
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return BusinessCategory{}, err
	}
	return s.Get(ctx, *id)
}
