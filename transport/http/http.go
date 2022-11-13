package transport

import (
	"context"
	"net/http"

	"github.com/eduardohoraciosanto/users/pkg/health"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type correlationID string

func NewHTTPRouter(hsvc health.Service) *mux.Router {

	hc := health.Handler{
		Service: hsvc,
	}

	r := mux.NewRouter()
	r.Use(correlationIDMiddleware)

	r.HandleFunc("/health", hc.Health).Methods(http.MethodGet)

	r.PathPrefix("/swagger").Handler(http.StripPrefix("/swagger", http.FileServer(http.Dir("./swagger"))))
	return r
}

func correlationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.Header.Get("X-Correlation-Id")
		if id == "" {
			// generate new version 4 uuid
			id = uuid.New().String()
		}
		// set the id to the request context
		ctx = context.WithValue(ctx, correlationID("correlation_id"), id)
		r = r.WithContext(ctx)

		// set the response header
		w.Header().Set("X-Correlation-Id", id)
		next.ServeHTTP(w, r)
	})
}
