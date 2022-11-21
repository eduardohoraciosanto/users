package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/eduardohoraciosanto/users/internal/config"
	"github.com/eduardohoraciosanto/users/internal/db"
	"github.com/eduardohoraciosanto/users/internal/logger"
	"github.com/eduardohoraciosanto/users/pkg/health"
	"github.com/eduardohoraciosanto/users/pkg/users"
	transport "github.com/eduardohoraciosanto/users/transport/http"
)

func main() {
	conf := config.New()

	l := logger.NewLogger("users api", config.GetVersion(), conf.TracingEnabled)

	/*
		cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
			// Hard coded credentials.
			awsConfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID: conf.AWSAccessKey, SecretAccessKey: conf.AWSSecretKey, SessionToken: "",
					Source: "environment",
				},
			}))
		if err != nil {
			l.WithError(err).Error(context.TODO(), "Unable to load AWS Config")
			return
		}
		db := db.NewDynamoDB("users", dynamodb.NewFromConfig(cfg), l.WithField("svc", "dynamo DB"))
	*/

	db := db.NewMemDB(l.WithField("svc", "db"))

	hsvc := health.NewService(
		l.WithField("svc", "health"),
		db,
	)

	userSvc := users.NewService(db, l.WithField("svc", "users"))

	httpTransportRouter := transport.NewHTTPRouter(hsvc, userSvc)

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", conf.Port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      httpTransportRouter,
	}

	l.WithField(
		"transport", "http").
		WithField(
			"port", conf.Port).
		Info(context.Background(), "Transport Start")

		// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			l.WithField(
				"transport", "http").
				WithError(err).
				Info(context.Background(), "Transport Stopped")
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	l.Info(context.Background(), "Service gracefully shut down")
	os.Exit(0)
}
