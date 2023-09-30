package business

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger) {
	res := resource{service, logger}
	r.HandleFunc("/api/v1/businesses/{id}", res.getByIdHandler).Methods("GET")
	//r.HandleFunc("/api/v1/businesses", res.getByNameHandler).Methods("GET")
	r.HandleFunc("/api/v1/businesses", res.create).Methods("POST")
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

//	func (r resource) getByNameHandler(w http.ResponseWriter, req *http.Request) {
//		name := req.URL.Query().Get("name")
//
//		category, _ := r.service.GetByName(req.Context(), name)
//		json.NewEncoder(w).Encode(category)
//	}
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
