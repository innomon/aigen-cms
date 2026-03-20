package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/formcms/formcms-go/core/api"
	"github.com/formcms/formcms-go/core/descriptors"
	"github.com/formcms/formcms-go/core/services"
	"github.com/formcms/formcms-go/infrastructure/relationdbdao"
)

func main() {
	// Initialize Database
	// For demo/testing, using SQLite. In production, this would come from config.
	dao, err := relationdbdao.CreateDao(descriptors.SQLite, "formcms.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dao.Close()

	// Initialize Services
	schemaService := services.NewSchemaService(dao)
	entityService := services.NewEntityService(schemaService, dao)
	graphqlService := services.NewGraphQLService(schemaService, entityService)

	// Initialize APIs
	schemaApi := api.NewSchemaApi(schemaService)
	entityApi := api.NewEntityApi(entityService)
	graphqlApi := api.NewGraphQLApi(graphqlService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("FormCMS Go is running!"))
	})

	// Register APIs
	schemaApi.Register(r)
	entityApi.Register(r)
	graphqlApi.Register(r)

	fmt.Println("Starting FormCMS Go on :5000...")
	log.Fatal(http.ListenAndServe(":5000", r))
}
