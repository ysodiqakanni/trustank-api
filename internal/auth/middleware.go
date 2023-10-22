package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	//routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ysodiqakanni/trustank-api/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func AuthenticateMiddleware(next http.Handler, jwtSecret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Your authentication logic here
		// Check if the JWT token is valid and extract user information
		// For example, you can check the "Authorization" header for the token
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization token missing", http.StatusUnauthorized)
			return
		}

		// test token validation
		var tokenSecret = []byte(jwtSecret)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method and return the secret key
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Invalid token signing method")
			}
			return tokenSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			fmt.Println(err.Error())
			return
		}
		// end test token validation. Let's inspect the token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		fmt.Println("User claims are: ", claims)
		// Add user information to the request context
		ctx := context.WithValue(r.Context(), "username", claims["name"].(string))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// User struct for context value
type User struct {
	ID       string
	Username string
	// Other user information
}

//****************************************

// Handler returns a JWT-based authentication middleware.
//func Handler(verificationKey string) routing.Handler {
//	return auth.JWT(verificationKey, auth.JWTOptions{TokenHandler: handleToken})
//}
//
//// handleToken stores the user identity in the request context so that it can be accessed elsewhere.
//func handleToken(c *routing.Context, token *jwt.Token) error {
//	ctx := WithUser(
//		c.Request.Context(),
//		token.Claims.(jwt.MapClaims)["id"].(primitive.ObjectID),
//		token.Claims.(jwt.MapClaims)["name"].(string),
//	)
//	c.Request = c.Request.WithContext(ctx)
//	return nil
//}

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
