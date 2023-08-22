package businessCategory

import (
	"context"
	"fmt"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/dbcontext"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository encapsulates the logic to access categories from the data source.
type Repository interface {
	// Get returns the category with the specified album ID.
	Get(ctx context.Context, id primitive.ObjectID) (entity.BusinessCategory, error)
	GetByName(ctx context.Context, id string) (entity.BusinessCategory, error)
}

// repository persists albums in database
type repository struct {
	collection *mongo.Collection
	logger     log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	col := db.DB().Collection("business_categories")
	logger.Infof("collection retrieved")
	return repository{col, logger}
}

func (r repository) Get(ctx context.Context, id primitive.ObjectID) (entity.BusinessCategory, error) {
	fmt.Println("Getting category by Id")
	filter := bson.M{"_id": id}
	var category entity.BusinessCategory
	err := r.collection.FindOne(ctx, filter).Decode(&category)

	return category, err
}

func (r repository) GetByName(ctx context.Context, name string) (entity.BusinessCategory, error) {
	fmt.Println("Getting category by name")
	filter := bson.M{"name": name}
	var category entity.BusinessCategory
	err := r.collection.FindOne(ctx, filter).Decode(&category)

	return category, err
}
