package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("FormCMS Go is running!"))
	})

	fmt.Println("Starting FormCMS Go on :5000...")
	log.Fatal(http.ListenAndServe(":5000", r))
}
