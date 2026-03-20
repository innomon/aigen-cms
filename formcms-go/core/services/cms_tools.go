package services

import (
	"context"
	"fmt"

	"github.com/formcms/formcms-go/utils/datamodels"
	"github.com/innomon/agentic/pkg/registry"
)

func RegisterCMSTools(entityService IEntityService, schemaService *SchemaService, a2uiService *A2UIService) {
	registry.RegisterToolHandler("cms_entity_list", func(ctx context.Context, args map[string]any) (any, error) {
		name, ok := args["name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid 'name' argument")
		}

		limit := "10"
		if l, ok := args["limit"].(string); ok {
			limit = l
		}

		records, _, err := entityService.List(ctx, name, datamodels.Pagination{Limit: &limit}, nil, nil)
		return records, err
	})

	registry.RegisterToolHandler("cms_entity_get", func(ctx context.Context, args map[string]any) (any, error) {
		name, ok := args["name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid 'name' argument")
		}
		id, ok := args["id"]
		if !ok {
			return nil, fmt.Errorf("missing 'id' argument")
		}

		record, err := entityService.Single(ctx, name, id)
		return record, err
	})

	registry.RegisterToolHandler("cms_entity_create", func(ctx context.Context, args map[string]any) (any, error) {
		name, ok := args["name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid 'name' argument")
		}
		data, ok := args["data"].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("missing or invalid 'data' argument")
		}

		record, err := entityService.Insert(ctx, name, data)
		return record, err
	})

	registry.RegisterToolHandler("cms_schema_list", func(ctx context.Context, args map[string]any) (any, error) {
		schemas, err := schemaService.All(ctx, nil, nil, nil)
		return schemas, err
	})

	registry.RegisterToolHandler("cms_a2ui_update", func(ctx context.Context, args map[string]any) (any, error) {
		id, ok := args["id"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid 'id' argument")
		}
		typ, ok := args["type"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid 'type' argument")
		}

		var attrs map[string]interface{}
		if a, ok := args["attributes"].(map[string]interface{}); ok {
			attrs = a
		} else if a, ok := args["attributes"].(map[string]any); ok {
			attrs = make(map[string]interface{})
			for k, v := range a {
				attrs[k] = v
			}
		}

		var children []string
		if c, ok := args["children"].([]interface{}); ok {
			for _, child := range c {
				if str, ok := child.(string); ok {
					children = append(children, str)
				}
			}
		} else if c, ok := args["children"].([]string); ok {
			children = c
		}

		comp := A2UIComponent{
			ID:         id,
			Type:       typ,
			Attributes: attrs,
			Children:   children,
		}

		a2uiService.UpdateComponent(ctx, comp)
		return "Component updated successfully", nil
	})
}
