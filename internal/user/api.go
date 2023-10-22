package user

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ysodiqakanni/trustank-api/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger) {
	res := resource{service, logger}
	r.HandleFunc("/api/v1/categories/{id}", res.getByIdHandler).Methods("GET")

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
