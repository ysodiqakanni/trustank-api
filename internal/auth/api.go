package auth

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"net/http"
)

type resource struct {
	service Service
	logger  log.Logger
}

// RegisterHandlers registers handlers for different HTTP requests.
//
//	func RegisterHandlers(rg *routing.RouteGroup, service Service, logger log.Logger) {
//		rg.Post("/login", login(service, logger))
//	}
func RegisterHandlers(r *mux.Router, service Service, logger log.Logger) {
	res := resource{service, logger}
	r.HandleFunc("/api/v1/login", res.loginHandler).Methods("POST")
}

// login returns a handler that handles user login request.
//func login(service Service, logger log.Logger) routing.Handler {
//	return func(c *routing.Context) error {
//		var req struct {
//			Username string `json:"username"`
//			Password string `json:"password"`
//		}
//
//		if err := c.Read(&req); err != nil {
//			logger.With(c.Request.Context()).Errorf("invalid request: %v", err)
//			return errors.BadRequest("")
//		}
//
//		token, err := service.Login(c.Request.Context(), req.Username, req.Password)
//		if err != nil {
//			return err
//		}
//		return c.Write(struct {
//			Token string `json:"token"`
//		}{token})
//	}
//}

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(token)
}
