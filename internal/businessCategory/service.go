package businessCategory

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/ysodiqakanni/trustank-api/internal/entity"
	"github.com/ysodiqakanni/trustank-api/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"time"
)

// Service encapsulates use case logic for businessCategories.
type Service interface {
	Get(ctx context.Context, id primitive.ObjectID) (BusinessCategory, error)
	GetByName(ctx context.Context, name string) (*BusinessCategory, error)
	Create(ctx context.Context, req CreateBusinessCategoryRequest) (BusinessCategory, error)
	Update(ctx context.Context, category UpdateBusinessCategoryRequest) (*entity.BusinessCategory, error)
	GetFeatured(ctx context.Context) []BusinessCategory
	Search(ctx context.Context, keyword string) []BusinessCategory
	Delete(ctx context.Context, categoryId string) error
}

// BusinessCategory represents the data about a BusinessCategory.
type BusinessCategory struct {
	entity.BusinessCategory
}

// CreateBusinessCategoryRequest represents an category creation request.
type CreateBusinessCategoryRequest struct {
	Name       string `json:"name"`
	IsFeatured bool   `json:"isFeatured"`
	IconUrl    string `json:"iconUrl"`
}
type UpdateBusinessCategoryRequest struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	IsFeatured bool   `json:"isFeatured"`
	IconUrl    string `json:"iconUrl"`
}

// Validate validates the CreateAlbumRequest fields.
func (m CreateBusinessCategoryRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128), validation.Match(regexp.MustCompile("^[a-zA-Z0-9].*$"))),
		validation.Field(&m.IconUrl, is.URL),
	)
}

func (m UpdateBusinessCategoryRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128), validation.Match(regexp.MustCompile("^[a-zA-Z0-9].*$"))),
		validation.Field(&m.IconUrl, is.URL),
		validation.Field(&m.Id, validation.Required),
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

func (s service) GetByName(ctx context.Context, name string) (*BusinessCategory, error) {
	category, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return &BusinessCategory{category}, nil
}
func (s service) Create(ctx context.Context, req CreateBusinessCategoryRequest) (BusinessCategory, error) {
	if err := req.Validate(); err != nil {
		return BusinessCategory{}, err
	}

	existing, _ := s.GetByName(ctx, req.Name)
	//emptyObj := BusinessCategory{}
	if existing != nil /*!= emptyObj*/ {
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
func (s service) Update(ctx context.Context, category UpdateBusinessCategoryRequest) (*entity.BusinessCategory, error) {
	objectId, _ := primitive.ObjectIDFromHex(category.Id)
	existingCategory, err := s.repo.Get(ctx, objectId)
	if err != nil {
		return nil, err
	}
	existingCategory.Name = category.Name
	existingCategory.IconUrl = category.IconUrl
	existingCategory.UpdatedAt = time.Now()
	fmt.Println("calling the repository layer for update")
	_, err = s.repo.Update(ctx, existingCategory)
	return &existingCategory, err
}

func (s service) Delete(ctx context.Context, categoryId string) error {
	objectId, _ := primitive.ObjectIDFromHex(categoryId)
	existingCategory, err := s.repo.Get(ctx, objectId)
	if err != nil {
		return err
	}
	existingCategory.IsDeleted = true
	fmt.Println("calling the repository layer for update")
	_, err = s.repo.Update(ctx, existingCategory)
	return err
}

func (s service) GetFeatured(ctx context.Context) []BusinessCategory {
	list := s.repo.GetFeaturedList(ctx)
	return list
}

func (s service) Search(ctx context.Context, keyword string) []BusinessCategory {
	list := s.repo.SearchCategories(ctx, keyword)
	return list
}
