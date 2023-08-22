package businessCategory

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service encapsulates use case logic for businessCategories.
type Service interface {
	Get(ctx context.Context, id primitive.ObjectID) (BusinessCategory, error)
	GetByName(ctx context.Context, name string) (BusinessCategory, error)
	//Query(ctx context.Context, offset, limit int) ([]Album, error)
	//Count(ctx context.Context) (int, error)
	//Create(ctx context.Context, input CreateAlbumRequest) (Album, error)
	//Update(ctx context.Context, id string, input UpdateAlbumRequest) (Album, error)
	//Delete(ctx context.Context, id string) (Album, error)
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
