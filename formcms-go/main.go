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
	"github.com/formcms/formcms-go/infrastructure/filestore"
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
	systemSettings := descriptors.DefaultSystemSettings()
	fileStore := filestore.NewLocalFileStore(systemSettings.LocalFileStoreOptions.PathPrefix, systemSettings.LocalFileStoreOptions.UrlPrefix)

	schemaService := services.NewSchemaService(dao)
	entityService := services.NewEntityService(schemaService, dao)
	graphqlService := services.NewGraphQLService(schemaService, entityService)
	assetService := services.NewAssetService(dao, fileStore, systemSettings)
	engagementService := services.NewEngagementService(dao)
	commentService := services.NewCommentService(dao)
	authService := services.NewAuthService(dao, "your-secret-key")
	notificationService := services.NewNotificationService(dao)
	auditService := services.NewAuditService(dao)
	pageService := services.NewPageService(schemaService, graphqlService)

	// Initialize APIs
	schemaApi := api.NewSchemaApi(schemaService)
	entityApi := api.NewEntityApi(entityService)
	graphqlApi := api.NewGraphQLApi(graphqlService)
	queryApi := api.NewQueryApi(graphqlService)
	assetApi := api.NewAssetApi(assetService)
	engagementApi := api.NewEngagementApi(engagementService)
	commentApi := api.NewCommentApi(commentService)
	authApi := api.NewAuthApi(authService)
	notificationApi := api.NewNotificationApi(notificationService, authApi)
	auditApi := api.NewAuditApi(auditService, authApi)
	staticApi := api.NewStaticApi()
	pageApi := api.NewPageApi(pageService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Register APIs
	schemaApi.Register(r)
	entityApi.Register(r)
	graphqlApi.Register(r)
	queryApi.Register(r)
	assetApi.Register(r)
	engagementApi.Register(r)
	commentApi.Register(r)
	authApi.Register(r)
	notificationApi.Register(r)
	auditApi.Register(r)
	staticApi.Register(r)
	pageApi.Register(r)

	fmt.Println("Starting FormCMS Go on :5000...")
	log.Fatal(http.ListenAndServe(":5000", r))
}
