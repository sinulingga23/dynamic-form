package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sinulingga23/dynamic-form/backend/db"
	"github.com/sinulingga23/dynamic-form/backend/implement/repository"
	"github.com/sinulingga23/dynamic-form/backend/implement/usecase"

	deliveryHttp "github.com/sinulingga23/dynamic-form/backend/delivery/http"
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

	// repository
	db, errConnect := db.ConnectDB()
	if errConnect != nil {
		log.Fatalf("errConnect: %v", errConnect)
	}
	pFormFieldRepository := repository.NewPFormFieldRepositoryImpl(db)

	// usecase
	pFormFieldUsecase := usecase.NewPFormFieldUsecase(db, pFormFieldRepository)

	// delivery - http
	formFieldHttp := deliveryHttp.NewFormFieldHttp(pFormFieldUsecase)
	formFieldHttp.ServeHandler(r)

	log.Printf("Running backend-service on: %s", port)
	if errListen := http.ListenAndServe(fmt.Sprintf(":%s", port), r); errListen != nil {
		log.Fatalf("Error when listen and serve backend-service on :%s, error: %v", port, errListen)
	}
}
