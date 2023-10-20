package business

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/qiangxue/go-rest-api/internal/auth"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

/*
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := verifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + err.Error()))
			return
		}
		name := claims.(jwt.MapClaims)["name"].(string)
		role := claims.(jwt.MapClaims)["role"].(string)

		r.Header.Set("name", name)
		r.Header.Set("role", role)

		next.ServeHTTP(w, r)
	})
}
*/

func PublicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Business open")
}

// ProtectedHandler is a protected handler.
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Business protected")
}
func RegisterBusinessHandlers(r *mux.Router, service Service, logger log.Logger, secret string) {
	r.HandleFunc("/v1/business1", PublicHandler).Methods("GET")

	// Protected Endpoint
	r.Handle("/v1/business2", auth.AuthenticateMiddleware(http.HandlerFunc(ProtectedHandler), secret)).Methods("GET")

}

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger, secret string) {
	res := resource{service, logger}

	r.HandleFunc("/api/v1/businesses/{id}", res.getByIdHandler).Methods("GET")

	// Protected Endpoint
	r.Handle("/api/v1/businesses", auth.AuthenticateMiddleware(http.HandlerFunc(res.getByNameHandler), secret)).Methods("GET")

	//
	//r.HandleFunc("/api/v1/businesses/{id}", res.getByIdHandler).Methods("GET")
	////r.Handle("/api/v1/businesses", authMiddleware(res.getByNameHandler)).Methods("GET")
	//
	//r.Use(authHandler111)
	//
	//r.HandleFunc("/api/v1/businesses", res.create).Methods("POST")

	//r.Use(func(next http.Handler) http.Handler {
	//	return authHandler(next)
	//})
	//r.Handle("/", authMiddleware(http.HandlerFunc(res.create)))
	//r.Handle("/api/levrai", authHandler(http.HandlerFunc(homeHandler)))
	//r.Handle("/api/levrai", authMiddleware(http.HandlerFunc(homeHandlerWrapper)))

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Matrix!"))
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getByIdHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	idk, _ := primitive.ObjectIDFromHex(id)

	business, _ := r.service.Get(req.Context(), idk)
	json.NewEncoder(w).Encode(business)
}

func (r resource) getByNameHandler(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")

	category, _ := r.service.GetByName(req.Context(), name)
	json.NewEncoder(w).Encode(category)
}
func (r resource) create(w http.ResponseWriter, req *http.Request) {
	var input CreateBusinessRequest

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	business, err := r.service.Register(req.Context(), input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(business)
}
