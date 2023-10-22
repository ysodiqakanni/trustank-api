package businessCategory

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ysodiqakanni/trustank-api/internal/auth"
	"github.com/ysodiqakanni/trustank-api/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger, secret string) {
	res := resource{service, logger}
	r.HandleFunc("/api/v1/categories/{id}", res.getByIdHandler).Methods("GET")
	r.HandleFunc("/api/v1/categories", res.getByNameHandler).Methods("GET")

	// Protected Endpoints
	r.Handle("/api/v1/categories", auth.AuthenticateMiddleware(http.HandlerFunc(res.create), secret)).Methods("POST")
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

	fmt.Println("calling the service layer")
	category, err := r.service.Create(req.Context(), input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
	//
	//categoryBytes, err := json.Marshal(category)
	//w.WriteHeader(http.StatusCreated)
	//_, err1 := w.Write(categoryBytes)
	//return err1
}
