package auth

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ysodiqakanni/trustank-api/pkg/log"
	"net/http"
)

type resource struct {
	service Service
	logger  log.Logger
}

// RegisterHandlers registers handlers for different HTTP requests.
func RegisterHandlers(r *mux.Router, service Service, logger log.Logger) {
	res := resource{service, logger}
	r.HandleFunc("/api/v1/login", res.loginHandler).Methods("POST")
}

func (r resource) loginHandler(w http.ResponseWriter, req *http.Request) {
	var input LoginRequest

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.logger.With(req.Context()).Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := input.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := r.service.Login(req.Context(), input.Username, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(token)
}
