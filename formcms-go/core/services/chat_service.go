package services

import (
	"context"
	"fmt"
	"log"

	"github.com/formcms/formcms-go/core/agentic/agents"
	"github.com/innomon/agentic/pkg/config"
	"github.com/innomon/agentic/pkg/registry"
)

type ChatService struct {
	Registry      *registry.Registry
	EntityService IEntityService
	SchemaService *SchemaService
	A2UIService   *A2UIService
}

func NewChatService(configPath string, entityService IEntityService, schemaService *SchemaService, a2uiService *A2UIService) (*ChatService, error) {
	// Register custom types and tools
	agents.RegisterRouterAgent()
	RegisterCMSTools(entityService, schemaService, a2uiService)

	// Load agentic config
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load agentic config: %w", err)
	}

	// Initialize registry
	reg := registry.New(cfg)

	svc := &ChatService{
		Registry:      reg,
		EntityService: entityService,
		SchemaService: schemaService,
		A2UIService:   a2uiService,
	}

	return svc, nil
}
