package transport

import (
	"context"
	"net/http"

	iContext "github.com/eduardohoraciosanto/users/internal/context"
	"github.com/eduardohoraciosanto/users/pkg/health"
	"github.com/eduardohoraciosanto/users/pkg/users"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func NewHTTPRouter(h health.Service, u users.Service) *mux.Router {

	hh := health.Handler{
		Service: h,
	}

	uh := users.Handler{
		Service: u,
	}

	r := mux.NewRouter()
	r.Use(correlationIDMiddleware)
	r.Use(remoteIPMiddleware)

	r.HandleFunc("/health", hh.Health).Methods(http.MethodGet)

	//User Related Endpoints
	r.HandleFunc("/users", uh.Create).Methods(http.MethodPost)
	r.HandleFunc("/users/login", uh.Login).Methods(http.MethodPost)

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
		ctx = context.WithValue(ctx, iContext.CorrelationID("correlation_id"), id)
		r = r.WithContext(ctx)

		// set the response header
		w.Header().Set("X-Correlation-Id", id)
		next.ServeHTTP(w, r)
	})
}

func remoteIPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ips := r.Header.Get("X-Forwarded-For")
		if ips != "" {
			ctx = context.WithValue(ctx, iContext.RemoteIP("remote_ip"), ips)
		} else {
			//not forwarded, lets get it from remote addr field
			ctx = context.WithValue(ctx, iContext.RemoteIP("remote_ip"), r.RemoteAddr)
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
