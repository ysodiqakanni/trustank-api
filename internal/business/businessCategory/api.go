package businessCategory

import (
	"fmt"
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/qiangxue/go-rest-api/pkg/log"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/categories/<id>", res.get)
	r.Get("/categories?name=<name>", res.getCategoryByName)

	r.Use(authHandler)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	id := c.Param("id")
	fmt.Println("using id: ", id)
	album, err := r.service.Get(c.Request.Context(), id)
	if err != nil {
		return err
	}

	return c.Write(album)
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
