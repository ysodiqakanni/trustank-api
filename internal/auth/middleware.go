package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/auth"
	"github.com/gorilla/mux"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/internal/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// region new middleware func
// Define a custom MiddlewareFunc type
type MiddlewareFunc func(http.Handler) http.Handler

func JwtMiddleware(next http.Handler, jwtSecret string) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")

			fmt.Println("Validating token")
			if tokenString == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Check the signing method and return the secret key
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Invalid token signing method")
				}
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// If the token is valid, proceed to the next handler
			handler.ServeHTTP(w, r)
		})
	}
}

// end region
func AuthenticateMiddleware(next http.Handler, jwtSecret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Your authentication logic here
		// Check if the JWT token is valid and extract user information
		fmt.Println("Checking auth")
		// For example, you can check the "Authorization" header for the token
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization token missing", http.StatusUnauthorized)
			return
		}

		// Validate the token and extract user information
		user, err := validateAndExtractToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// test token validation
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method and return the secret key
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Invalid token signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		// end test token validation

		// Add user information to the request context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Example function to validate and extract token
func validateAndExtractToken(tokenString string) (*User, error) {
	// Your JWT validation logic here
	// Parse and validate the token, extract user details
	usr := User{}
	return &usr, nil
}

// User struct for context value
type User struct {
	ID       string
	Username string
	// Other user information
}

//****************************************

// Handler returns a JWT-based authentication middleware.
func Handler(verificationKey string) routing.Handler {
	return auth.JWT(verificationKey, auth.JWTOptions{TokenHandler: handleToken})
}

// handleToken stores the user identity in the request context so that it can be accessed elsewhere.
func handleToken(c *routing.Context, token *jwt.Token) error {
	ctx := WithUser(
		c.Request.Context(),
		token.Claims.(jwt.MapClaims)["id"].(primitive.ObjectID),
		token.Claims.(jwt.MapClaims)["name"].(string),
	)
	c.Request = c.Request.WithContext(ctx)
	return nil
}

type contextKey int

const (
	userKey contextKey = iota
)

// WithUser returns a context that contains the user identity from the given JWT.
func WithUser(ctx context.Context, id primitive.ObjectID, name string) context.Context {
	return context.WithValue(ctx, userKey, entity.User{ID: id, Name: name})
}

// CurrentUser returns the user identity from the given context.
// Nil is returned if no user identity is found in the context.
func CurrentUser(ctx context.Context) Identity {
	if user, ok := ctx.Value(userKey).(entity.User); ok {
		return user
	}
	return nil
}

// MockAuthHandler creates a mock authentication middleware for testing purpose.
// If the request contains an Authorization header whose value is "TEST", then
// it considers the user is authenticated as "Tester" whose ID is "100".
// It fails the authentication otherwise.
func MockAuthHandler(c *routing.Context) error {
	if c.Request.Header.Get("Authorization") != "TEST" {
		return errors.Unauthorized("")
	}
	ctx := WithUser(c.Request.Context(), primitive.NewObjectID(), "Tester")
	c.Request = c.Request.WithContext(ctx)
	return nil
}

// MockAuthHeader returns an HTTP header that can pass the authentication check by MockAuthHandler.
func MockAuthHeader() http.Header {
	header := http.Header{}
	header.Add("Authorization", "TEST")
	return header
}
