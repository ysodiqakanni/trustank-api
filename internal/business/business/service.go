package business

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/internal/user"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"golang.org/x/crypto/bcrypt"
)

// Service encapsulates use case logic for businesses.
type Service interface {
	Get(ctx context.Context, id primitive.ObjectID) (Business, error)
	Register(ctx context.Context, req CreateBusinessRequest) (Business, error)
}

// Business represents the data about a BusinessCategory.
type Business struct {
	entity.Business
}

// CreateBusinessCategoryRequest represents an category creation request.
type CreateBusinessRequest struct {
	BusinessName    string `json:"businessName,omitempty" validate:"required"`
	Website         string `json:"website,omitempty" validate:"url"`
	OwnerFullName   string `json:"ownerFullName,omitempty" validate:"required"`
	OwnerJobTitle   string
	WorkEmail       string `json:"workEmail,omitempty" validate:"required,email"`
	PhoneNumber     string `json:"phoneNumber,omitempty" validate:"required"`
	Password        string `json:"password,omitempty" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword,omitempty" validate:"required,eqfield=Password"`
}

// Validate validates the CreateAlbumRequest fields.
func (m CreateBusinessRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.BusinessName, validation.Required, validation.Length(2, 128)),
		validation.Field(&m.OwnerFullName, validation.Required, validation.Length(5, 128)),
		validation.Field(&m.WorkEmail, validation.Required, validation.Length(7, 128), is.Email),
		validation.Field(&m.Password, validation.Required, validation.Length(4, 128)),

		//validation.Field(&a.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
	)
}

type service struct {
	repo     Repository
	userRepo user.Repository
	logger   log.Logger
}

// NewService creates a new category service.
func NewService(repo Repository, userRepo user.Repository, logger log.Logger) Service {
	return service{repo, userRepo, logger}
}

// Get returns the album with the specified the album ID.
func (s service) Get(ctx context.Context, id primitive.ObjectID) (Business, error) {
	business, err := s.repo.Get(ctx, id)
	if err != nil {
		return Business{}, err
	}
	return Business{business}, nil
}

func (s service) Register(ctx context.Context, req CreateBusinessRequest) (Business, error) {
	if err := req.Validate(); err != nil {
		return Business{}, err
	}
	fmt.Println("After validation with email: ", req.WorkEmail, primitive.ObjectID{})

	// check if a user with that name exists
	existing, err := s.userRepo.GetByEmail(ctx, req.WorkEmail)
	fmt.Println("get by email: ", existing, err)
	emptyId := primitive.ObjectID{}
	if err == nil || existing.ID != emptyId {
		return Business{}, errors.New("A business with this email already exists")
	}

	//emptyUserObject := entity.User{}
	//if existing != emptyUserObject {
	//	http.Error(w, "A user with this email already exists", http.StatusConflict)
	//	return
	//}
	// now send the data to the repo
	fmt.Println("After uniqueness check")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return Business{}, err
	}

	transactionOptions := options.Transaction().
		SetReadConcern(readconcern.Majority()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	// Start a new session
	session, err := s.repo.StartSession()
	if err != nil {
		s.logger.Error(err)
	}
	defer session.EndSession(context.Background())

	// Start the transaction
	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		// Start the transaction
		err := session.StartTransaction(transactionOptions)
		if err != nil {
			return err
		}

		// Create a user
		// let's create a user object from this
		user := entity.User{
			Email:          req.WorkEmail,
			Name:           req.OwnerFullName,
			HashedPassword: hashedPassword,
			Role:           []string{"business"},
		}
		// Insert the user document
		_, err = s.userRepo.Create(sessionContext, user)
		if err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}

		// Create a business object
		business := entity.Business{
			Name:          req.BusinessName,
			Email:         req.WorkEmail,
			Website:       req.Website,
			OwnerId:       user.ID,
			OwnerName:     req.OwnerFullName,
			OwnerJobTitle: req.OwnerJobTitle,
		}

		// Insert the profile document
		_, err = s.repo.Create(sessionContext, business)
		if err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}

		// Commit the transaction
		err = session.CommitTransaction(sessionContext)
		if err != nil {
			return err
		}

		return nil
	})

	return Business{}, nil
}
