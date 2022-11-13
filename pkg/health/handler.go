package health

import (
	"net/http"

	"github.com/eduardohoraciosanto/users/internal/response"
)

type Handler struct {
	Service Service
}

//Health is the handler for the health endpoint
func (c *Handler) Health(w http.ResponseWriter, r *http.Request) {
	//using lower level pkg to do the logic
	service, err := c.Service.HealthCheck(r.Context())
	if err != nil {
		response.RespondWithError(w, response.StandardInternalServerError)
		return
	}
	hr := HealthResponse{
		Services: []TransportHealth{
			{
				Name:  "service",
				Alive: service,
			},
		},
	}
	response.RespondWithData(w, http.StatusOK, hr)
}
