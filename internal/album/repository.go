package album

import (
	"context"
	"fmt"
	"github.com/ysodiqakanni/trustank-api/internal/entity"
	"github.com/ysodiqakanni/trustank-api/pkg/dbcontext"
	"github.com/ysodiqakanni/trustank-api/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id string) (entity.Album, error)
	// Count returns the number of albums.
	//Count(ctx context.Context) (int, error)
	//// Query returns the list of albums with the given offset and limit.
	//Query(ctx context.Context, offset, limit int) ([]entity.Album, error)
	// Create saves a new album in the storage.
	Create(ctx context.Context, album entity.Album) error
	// Update updates the album with given ID in the storage.
	Update(ctx context.Context, album entity.Album) error
	// Delete removes the album with given ID from the storage.
	Delete(ctx context.Context, id string) error
}

// repository persists albums in database
type repository struct {
	db     *dbcontext.DB
	colln  *mongo.Collection
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	col := db.DB().Collection("albums")
	logger.Infof("collection retrieved")
	return repository{db, col, logger}
}

func (r repository) Get(ctx context.Context, id string) (entity.Album, error) {
	fmt.Println("In get repo")
	filter := bson.M{"_id": id}
	var album entity.Album
	err := r.colln.FindOne(ctx, filter).Decode(&album)

	return album, err
}

// NewRepository creates a new album repository
//func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
//	return repository{db, logger}
//}
//
//// Get reads the album with the specified ID from the database.
//func (r repository) Get(ctx context.Context, id string) (entity.Album, error) {
//	var album entity.Album
//	err := r.db.With(ctx).Select().Model(id, &album)
//	return album, err
//}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r repository) Create(ctx context.Context, album entity.Album) error {
	_, err := r.colln.InsertOne(ctx, album)
	return err
	// refactor to return both error and Id of the new data

	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("inserted document with ID %v\n", res.InsertedID)
	//return r.db.With(ctx).Model(&album).Insert()
}

// Update saves the changes to an album in the database.
func (r repository) Update(ctx context.Context, album entity.Album) error {
	return nil
	//return r.db.With(ctx).Model(&album).Update()
}

// // Delete deletes an album with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id string) error {
	return nil
	//album, err := r.Get(ctx, id)
	//if err != nil {
	//	return err
	//}
	//return r.db.With(ctx).Model(&album).Delete()
}

// Count returns the number of the album records in the database.
//func (r repository) Count(ctx context.Context) (int, error) {
//	var count int
//	err := r.db.With(ctx).Select("COUNT(*)").From("album").Row(&count)
//	return count, err
//}
//
//// Query retrieves the album records with the specified offset and limit from the database.
//func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.Album, error) {
//	var albums []entity.Album
//	err := r.db.With(ctx).
//		Select().
//		OrderBy("id").
//		Offset(int64(offset)).
//		Limit(int64(limit)).
//		All(&albums)
//	return albums, err
//}
