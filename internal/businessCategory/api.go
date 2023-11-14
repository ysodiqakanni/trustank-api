package businessCategory

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/qiangxue/go-rest-api/internal/auth"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger, secret string) {
	res := resource{service, logger}
	//r.HandleFunc("/api/v1/categories/{id}", res.getByIdHandler).Methods("GET")
	r.HandleFunc("/api/v1/categories", res.getByNameHandler).Methods("GET")
	r.HandleFunc("/api/v1/categories/search", res.searchCategoriesHandler).Methods("GET")
	r.HandleFunc("/api/v1/categories/featured", res.getFeaturedCategoriesHandler).Methods("GET")

	// Protected Endpoints
	//r.Handle("/api/v1/categories", auth.AuthenticateMiddleware(http.HandlerFunc(res.create), secret)).Methods("POST")
	r.Handle("/api/v1/categories", auth.AuthenticateMiddleware(auth.RoleMiddleware(http.HandlerFunc(res.create), "admin"), secret)).Methods("POST")
	r.Handle("/api/v1/categories", auth.AuthenticateMiddleware(http.HandlerFunc(res.updateCategoryHandler), secret)).Methods("PUT")
	r.Handle("/api/v1/categories", auth.AuthenticateMiddleware(http.HandlerFunc(res.deleteCategoryHandler), secret)).Methods("DELETE")
	r.Use()
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getByIdHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	idk, _ := primitive.ObjectIDFromHex(id)

	category, _ := r.service.Get(req.Context(), idk)
	json.NewEncoder(w).Encode(category)
}

func (r resource) getByNameHandler(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")

	category, _ := r.service.GetByName(req.Context(), name)
	json.NewEncoder(w).Encode(category)
}
func (r resource) create(w http.ResponseWriter, req *http.Request) {
	var input CreateBusinessCategoryRequest

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("calling the service layer")
	category, err := r.service.Create(req.Context(), input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (r resource) updateCategoryHandler(w http.ResponseWriter, req *http.Request) {
	// update category data
	var input UpdateBusinessCategoryRequest

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	test, err := r.service.Update(req.Context(), input)
	if err != nil {
		r.logger.Error(err)
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(test)
}
func (r resource) deleteCategoryHandler(w http.ResponseWriter, req *http.Request) {
	// delete a category by Id
	// first get the category by Id
	vars := mux.Vars(req)
	id := vars["id"]
	err := r.service.Delete(req.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// search categories by query. Paginated
func (r resource) searchCategoriesHandler(w http.ResponseWriter, req *http.Request) {
	// for now we use a simple contain to search categories
	//vars := mux.Vars(req)
	//keyword := vars["query"]
	keyword := req.URL.Query().Get("query")
	results := r.service.Search(req.Context(), keyword)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

// get list of featured categories. Paginate?
func (r resource) getFeaturedCategoriesHandler(w http.ResponseWriter, req *http.Request) {
	results := r.service.GetFeatured(req.Context())
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}
