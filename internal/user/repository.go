package user

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
	Get(ctx context.Context, id primitive.ObjectID) (entity.User, error)
	GetByEmail(ctx context.Context, id string) (entity.User, error)
	Create(ctx context.Context, user entity.User) (*primitive.ObjectID, error)
	StartSession() (mongo.Session, error)
}

// repository persists albums in database
type repository struct {
	collection *mongo.Collection
	logger     log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	col := db.DB().Collection("users")
	return repository{col, logger}
}

func (r repository) StartSession() (mongo.Session, error) {
	return r.collection.Database().Client().StartSession()
}

func (r repository) Get(ctx context.Context, id primitive.ObjectID) (entity.User, error) {
	filter := bson.M{"_id": id}
	var user entity.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)

	return user, err
}
func (r repository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	filter := bson.M{"email": bson.M{"$regex": primitive.Regex{Pattern: "^" + email + "$", Options: "i"}}}
	var user entity.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)

	fmt.Println("user data: ", user)
	return user, err
}
func (r repository) Create(ctx context.Context, user entity.User) (*primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	fmt.Printf("inserted user data with ID %v\n", result.InsertedID)
	id := result.InsertedID.(primitive.ObjectID)
	return &id, err
}
