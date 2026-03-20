package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/formcms/formcms-go/core/api"
	"github.com/formcms/formcms-go/core/apps"
	"github.com/formcms/formcms-go/core/descriptors"
	"github.com/formcms/formcms-go/core/services"
	"github.com/formcms/formcms-go/infrastructure/filestore"
	"github.com/formcms/formcms-go/infrastructure/relationdbdao"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/acme/autocert"
)

func isExternalDomain(domain string) bool {
	if domain == "" || domain == "localhost" {
		return false
	}
	// Check if it's an IP address
	if net.ParseIP(domain) != nil {
		return false
	}
	return true
}

func main() {
	// Initialize Database
	// For demo/testing, using SQLite. In production, this would come from config.
	dao, err := relationdbdao.CreateDao(descriptors.SQLite, "formcms.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dao.Close()

	// Ensure core tables exist
	_, err = dao.GetDb().ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS __schemas (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			schema_id TEXT,
			name TEXT,
			type TEXT,
			settings TEXT,
			description TEXT,
			is_latest BOOLEAN,
			publication_status TEXT,
			created_at DATETIME,
			created_by TEXT,
			deleted BOOLEAN
		);
		CREATE TABLE IF NOT EXISTS __users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			avatar_path TEXT,
			created_at DATETIME,
			updated_at DATETIME
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Services
	systemSettings := descriptors.DefaultSystemSettings()
	fileStore := filestore.NewLocalFileStore(systemSettings.LocalFileStoreOptions.PathPrefix, systemSettings.LocalFileStoreOptions.UrlPrefix)

	schemaService := services.NewSchemaService(dao)

	enabledApps, err := apps.LoadAppsConfig()
	if err != nil {
		log.Fatalf("Failed to load apps config: %v", err)
	}

	for _, appName := range enabledApps {
		log.Printf("Setting up app schemas: %s", appName)
		if err := apps.SetupApp(context.Background(), appName, schemaService, dao); err != nil {
			log.Printf("Warning: failed to setup app %s schemas: %v\n", appName, err)
		}
	}

	entityService := services.NewEntityService(schemaService, dao)

	for _, appName := range enabledApps {
		log.Printf("Setting up test data for app: %s", appName)
		if err := apps.SetupAppTestData(context.Background(), appName, entityService); err != nil {
			log.Printf("Warning: failed to setup test data for app %s: %v\n", appName, err)
		}
	}

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

	domain := os.Getenv("DOMAIN")
	if isExternalDomain(domain) {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domain),
			Cache:      autocert.DirCache("certs"),
		}

		server := &http.Server{
			Addr:    ":443",
			Handler: r,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		fmt.Printf("Starting FormCMS Go on %s with autocert...\n", domain)
		// Redirect HTTP to HTTPS
		go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		port := os.Getenv("PORT")
		if port == "" {
			port = "5000"
		}
		fmt.Printf("Starting FormCMS Go on :%s...\n", port)
		log.Fatal(http.ListenAndServe(":"+port, r))
	}
}
