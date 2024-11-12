// server.go
//go:build !test
// +build !test

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Black-And-White-Club/tcr-bot-user-service/graph"
	"github.com/Black-And-White-Club/tcr-bot-user-service/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	// Import pgxpool for PostgreSQL connection pooling
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize PostgreSQL client
	ctx := context.Background()
	dataSourceName := os.Getenv("DATABASE_URL") // Ensure this environment variable is set

	// Create the PostgreSQL client
	pgClient, err := service.NewPGClient(dataSourceName)
	if err != nil {
		log.Fatalf("Failed to create PostgreSQL client: %v", err)
	}
	defer pgClient.Close(ctx) // Pass context here

	// Create UserService
	userService := service.NewUserService(pgClient) // Assume you have a UserService struct

	// Create a new Chi router
	router := chi.NewRouter()

	// Use Chi's built-in middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Create a new GraphQL server with the resolver that has the UserService
	gqlServer := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		UserService: userService,
	}}))

	// Set up routes
	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.Handle("/graphql", gqlServer)

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Start the HTTP server
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}
