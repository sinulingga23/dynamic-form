package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	port = "8085"
)

func init() {
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("it's work!"))
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Running backend-service on: %s", port)
	if errListen := http.ListenAndServe(fmt.Sprintf(":%s", port), r); errListen != nil {
		log.Fatalf("Error when listen and serve backend-service on :%s, error: %v", port, errListen)
	}
}
