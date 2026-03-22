package framework

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/innomon/aigen-cms/core/api"
	"github.com/innomon/aigen-cms/core/apps"
	"github.com/innomon/aigen-cms/core/descriptors"
	"github.com/innomon/aigen-cms/core/services"
	"github.com/innomon/aigen-cms/infrastructure/filestore"
	"github.com/innomon/aigen-cms/infrastructure/relationdbdao"
	chi "github.com/go-chi/chi/v5"
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

func Start(cfg *Config) error {
	// Initialize Database
	dao, err := relationdbdao.CreateDao(descriptors.DatabaseProvider(cfg.DatabaseType), cfg.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("failed to create dao: %w", err)
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
		CREATE TABLE IF NOT EXISTS __user_channels (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			channel_type TEXT NOT NULL,
			identifier TEXT NOT NULL,
			is_authenticated BOOLEAN DEFAULT 0,
			metadata TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES __users(id)
		);
		CREATE TABLE IF NOT EXISTS __auth_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			channel_type TEXT NOT NULL,
			action TEXT NOT NULL,
			ip_address TEXT,
			user_agent TEXT,
			success BOOLEAN,
			metadata TEXT,
			created_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES __users(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create core tables: %w", err)
	}

	// Initialize Services
	systemSettings := descriptors.DefaultSystemSettings()
	systemSettings.LocalFileStoreOptions.PathPrefix = filepath.Join(cfg.WWWRoot, "files")
	fileStore := filestore.NewLocalFileStore(systemSettings.LocalFileStoreOptions.PathPrefix, systemSettings.LocalFileStoreOptions.UrlPrefix)

	schemaService := services.NewSchemaService(dao)
	permissionService := services.NewPermissionService(dao, schemaService)

	enabledApps, err := apps.LoadAppsConfig(cfg.AppsDir)
	if err != nil {
		log.Printf("Warning: failed to load apps config from %s: %v", cfg.AppsDir, err)
		enabledApps = []string{} // Proceed without apps if failed
	}

	for _, appName := range enabledApps {
		log.Printf("Setting up app schemas: %s", appName)
		if err := apps.SetupApp(context.Background(), cfg.AppsDir, appName, schemaService, dao); err != nil {
			log.Printf("Warning: failed to setup app %s schemas: %v\n", appName, err)
		}
	}

	entityService := services.NewEntityService(schemaService, dao, permissionService)

	for _, appName := range enabledApps {
		log.Printf("Setting up test data for app: %s", appName)
		if err := apps.SetupAppTestData(context.Background(), cfg.AppsDir, appName, entityService); err != nil {
			log.Printf("Warning: failed to setup test data for app %s: %v\n", appName, err)
		}
	}

	graphqlService := services.NewGraphQLService(schemaService, entityService)
	assetService := services.NewAssetService(dao, fileStore, systemSettings)
	engagementService := services.NewEngagementService(dao)
	commentService := services.NewCommentService(dao)
	notificationService := services.NewNotificationService(dao)
	channelService := services.NewChannelService(dao, cfg.Channels)
	authService := services.NewAuthService(dao, "your-secret-key", channelService)
	auditService := services.NewAuditService(dao)
	pageService := services.NewPageService(schemaService, graphqlService)
	a2uiService := services.NewA2UIService()

	chatService, err := services.NewChatService(cfg.AgenticConfigPath, entityService, schemaService, a2uiService)
	if err != nil {
		log.Printf("Warning: failed to initialize chat service (agentic config missing or invalid): %v", err)
	}

	// Initialize Prototype Components for A2UI
	a2uiService.UpdateComponent(context.Background(), services.A2UIComponent{
		ID:   "root",
		Type: "Column",
		Children: []string{"header", "card-1", "card-2", "card-3"},
	})
	a2uiService.UpdateComponent(context.Background(), services.A2UIComponent{
		ID:   "header",
		Type: "Heading",
		Attributes: map[string]interface{}{
			"content": "A2UI Agent Intelligence Hub",
			"level":   1,
		},
	})
	a2uiService.UpdateComponent(context.Background(), services.A2UIComponent{
		ID:   "card-1",
		Type: "Card",
		Children: []string{"counter-text", "increment-btn"},
	})
	a2uiService.UpdateComponent(context.Background(), services.A2UIComponent{
		ID:   "counter-text",
		Type: "Text",
		Attributes: map[string]interface{}{
			"content": "Live System Counter: 0",
			"count":   0.0,
		},
	})
	a2uiService.UpdateComponent(context.Background(), services.A2UIComponent{
		ID:   "increment-btn",
		Type: "Button",
		Attributes: map[string]interface{}{
			"label":   "Trigger Agent Signal",
			"variant": "success",
			"action":  "increment",
		},
	})

	// Add Data Table Card
	a2uiService.UpdateComponent(context.Background(), services.A2UIComponent{
		ID:   "card-2",
		Type: "Card",
		Children: []string{"table-title", "audit-table"},
	})
	a2uiService.UpdateComponent(context.Background(), services.A2UIComponent{
		ID:   "table-title",
		Type: "Heading",
		Attributes: map[string]interface{}{
			"content": "Recent System Activity",
			"level":   4,
		},
	})
	a2uiService.UpdateComponent(context.Background(), services.A2UIComponent{
		ID:   "audit-table",
		Type: "DataTable",
		Attributes: map[string]interface{}{
			"columns": []string{"User", "Action", "Timestamp"},
			"rows": []map[string]interface{}{
				{"User": "admin@example.com", "Action": "Login", "Timestamp": "2024-03-20 10:00:01"},
				{"User": "editor@example.com", "Action": "Update Lead", "Timestamp": "2024-03-20 10:15:22"},
				{"User": "admin@example.com", "Action": "Delete Log", "Timestamp": "2024-03-20 11:30:00"},
			},
		},
	})

	// Add Graph Card
	a2uiService.UpdateComponent(context.Background(), services.A2UIComponent{
		ID:   "card-3",
		Type: "Card",
		Children: []string{"chart-title", "traffic-chart"},
	})
	a2uiService.UpdateComponent(context.Background(), services.A2UIComponent{
		ID:   "chart-title",
		Type: "Heading",
		Attributes: map[string]interface{}{
			"content": "Traffic Analytics",
			"level":   4,
		},
	})
	a2uiService.UpdateComponent(context.Background(), services.A2UIComponent{
		ID:   "traffic-chart",
		Type: "Chart",
		Attributes: map[string]interface{}{
			"chartType": "line",
			"label":     "Requests per Minute",
			"labels":    []string{"10:00", "10:10", "10:20", "10:30", "10:40", "10:50"},
			"data":      []float64{12, 19, 3, 5, 2, 3},
		},
	})

	// Initialize APIs
	authApi := api.NewAuthApi(authService, permissionService)
	rbacApi := api.NewRBACApi(entityService, authApi)
	schemaApi := api.NewSchemaApi(schemaService, authApi)
	entityApi := api.NewEntityApi(entityService, authApi)
	graphqlApi := api.NewGraphQLApi(graphqlService, authApi)
	queryApi := api.NewQueryApi(graphqlService, authApi)
	assetApi := api.NewAssetApi(assetService)
	engagementApi := api.NewEngagementApi(engagementService, authApi)
	commentApi := api.NewCommentApi(commentService, authApi)
	notificationApi := api.NewNotificationApi(notificationService, authApi)
	auditApi := api.NewAuditApi(auditService, authApi)
	channelApi := api.NewChannelApi(channelService, authApi)
	staticApi := api.NewStaticApi()
	pageApi := api.NewPageApi(pageService, authService, authApi)
	a2uiApi := api.NewA2UIApi(a2uiService, authApi)
	var chatApi *api.ChatApi
	if chatService != nil {
		chatApi = api.NewChatApi(chatService, authApi)
	}

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
	rbacApi.Register(r)
	notificationApi.Register(r)
	auditApi.Register(r)
	channelApi.Register(r)
	staticApi.Register(r)
	pageApi.Register(r)
	a2uiApi.Register(r)
	if chatApi != nil {
		chatApi.Register(r)
	}

	if isExternalDomain(cfg.Domain) {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.Domain),
			Cache:      autocert.DirCache("certs"),
		}

		server := &http.Server{
			Addr:    ":443",
			Handler: r,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		fmt.Printf("Starting AiGen CMS on %s with autocert...\n", cfg.Domain)
		// Redirect HTTP to HTTPS
		go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
		return server.ListenAndServeTLS("", "")
	} else {
		fmt.Printf("Starting AiGen CMS on :%s...\n", cfg.Port)
		return http.ListenAndServe(":"+cfg.Port, r)
	}
}
