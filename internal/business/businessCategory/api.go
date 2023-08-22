package businessCategory

import (
	"encoding/json"
	"fmt"
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/gorilla/mux"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/categories/<id>", res.get)
	r.Get("/categories?name=<name>", res.getCategoryByName)

	// for routes requiring auth
	r.Use(authHandler)
}

func RegisterHandlersMux(r *mux.Router, service Service, logger log.Logger) {
	res := resource{service, logger}
	r.HandleFunc("/api/v1/categories/{id}", res.GetByIdHandler).Methods("GET")
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	id := c.Param("id")
	fmt.Println("using id: ", id)
	album, err := r.service.Get(c.Request.Context(), primitive.NewObjectID())
	if err != nil {
		return err
	}

	return c.Write(album)
}

func (r resource) GetByIdHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("In byId handler")
	vars := mux.Vars(req)
	id := vars["id"]
	fmt.Printf("retrieving category with id: %v \n", id)

	idk, _ := primitive.ObjectIDFromHex(id)
	fmt.Println("Id in string: ", idk.Hex())

	category, _ := r.service.Get(req.Context(), idk)
	//result := res.service.Get().catRepo.GetBusinessCategoryById(idk)
	json.NewEncoder(w).Encode(category)
}

func (r resource) getCategoryByName(c *routing.Context) error {
	name := c.Param("name")
	fmt.Println("using name: ", name)
	album, err := r.service.GetByName(c.Request.Context(), name)
	if err != nil {
		return err
	}

	return c.Write(album)
}
