package business

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
	Get(ctx context.Context, id primitive.ObjectID) (entity.Business, error)
	GetByEmail(ctx context.Context, email string) (entity.Business, error)
	Create(ctx context.Context, business entity.Business) (*primitive.ObjectID, error)
	StartSession() (mongo.Session, error)
}

// repository persists albums in database
type repository struct {
	collection *mongo.Collection
	logger     log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	col := db.DB().Collection("businesses")
	return repository{col, logger}
}

func (r repository) StartSession() (mongo.Session, error) {
	return r.collection.Database().Client().StartSession()
}

func (r repository) Get(ctx context.Context, id primitive.ObjectID) (entity.Business, error) {
	filter := bson.M{"_id": id}
	var business entity.Business
	err := r.collection.FindOne(ctx, filter).Decode(&business)

	return business, err
}

func (r repository) GetByEmail(ctx context.Context, email string) (entity.Business, error) {
	filter := bson.M{"email": bson.M{"$regex": primitive.Regex{Pattern: "^" + email + "$", Options: "i"}}}
	var business entity.Business
	err := r.collection.FindOne(ctx, filter).Decode(&business)

	return business, err
}

func (r repository) Create(ctx context.Context, category entity.Business) (*primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}

	fmt.Printf("inserted document with ID %v\n", result.InsertedID)
	id := result.InsertedID.(primitive.ObjectID)
	return &id, err
}
