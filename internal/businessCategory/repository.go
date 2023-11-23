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
	Create(ctx context.Context, category entity.BusinessCategory) (*primitive.ObjectID, error)
	Update(ctx context.Context, category entity.BusinessCategory) (*primitive.ObjectID, error)
	GetFeaturedList(ctx context.Context) []BusinessCategory
	SearchCategories(ctx context.Context, keyword string) []BusinessCategory

	StartSession() (mongo.Session, error)
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

func (r repository) StartSession() (mongo.Session, error) {
	return r.collection.Database().Client().StartSession()
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
	filter := bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: "^" + name + "$", Options: "i"}}}
	var category entity.BusinessCategory
	err := r.collection.FindOne(ctx, filter).Decode(&category)

	return category, err
}
func (r repository) Create(ctx context.Context, category entity.BusinessCategory) (*primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}

	fmt.Printf("inserted document with ID %v\n", result.InsertedID)
	id := result.InsertedID.(primitive.ObjectID)
	return &id, err
}

func (r repository) Update(ctx context.Context, category entity.BusinessCategory) (*primitive.ObjectID, error) {
	filter := bson.M{"_id": category.ID}
	updateDoc, err := bson.Marshal(category)
	if err != nil {
		// Todo: log error
		return nil, err
	}
	result, err := r.collection.ReplaceOne(context.TODO(), filter, updateDoc)
	if err != nil {
		return nil, err
	}

	fmt.Printf("document updated")
	id := result.UpsertedID.(primitive.ObjectID)
	return &id, err
}

func (r repository) GetFeaturedList(ctx context.Context) []BusinessCategory {
	filter := bson.D{}
	cursor, err := r.collection.Find(ctx, filter)

	if err != nil {
		r.logger.Error(err)
		return []BusinessCategory{}
	}
	defer cursor.Close(ctx)

	categories := CursorToBusinessCategories(ctx, cursor)
	return categories
}
func (r repository) SearchCategories(ctx context.Context, keyword string) []BusinessCategory {
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": keyword, "$options": "i"}},
			{"description": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error(err)
		return []BusinessCategory{}
	}
	defer cursor.Close(ctx)

	categories := CursorToBusinessCategories(ctx, cursor)
	return categories
}

func CursorToBusinessCategories(ctx context.Context, cursor *mongo.Cursor) []BusinessCategory {
	var categories []BusinessCategory = []BusinessCategory{}

	// Iterate through the cursor and unmarshal each document into a BusinessCategory struct
	for cursor.Next(ctx) {
		var category entity.BusinessCategory
		err := cursor.Decode(&category)
		if err != nil {
			fmt.Errorf(err.Error())
			return []BusinessCategory{}
		}
		categories = append(categories, BusinessCategory{category})
	}

	if err := cursor.Err(); err != nil {
		// Handle any cursor error
		fmt.Errorf(err.Error())
		return []BusinessCategory{}
	}
	return categories
}
