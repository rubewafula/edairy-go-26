package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	initializers "github.com/rubewafula/edairy-go-26/internal/initializers"
	"github.com/rubewafula/edairy-go-26/internal/routes"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.InitTimezone()
	db.ConnectToDatabase()
}

func main() {
	router := routes.SetupRouter()

	addr := "0.0.0.0:" + initializers.GetEnv("PORT", "8000")
	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown: %v", err)
	}

	if err := db.CloseDatabase(); err != nil {
		log.Printf("database close: %v", err)
	}

	log.Println("server stopped")
}
