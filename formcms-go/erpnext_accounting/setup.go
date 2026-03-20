package erpnext_accounting

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/formcms/formcms-go/core/descriptors"
	"github.com/formcms/formcms-go/core/services"
	"github.com/formcms/formcms-go/infrastructure/relationdbdao"
	"github.com/formcms/formcms-go/utils/datamodels"
)

func Setup(ctx context.Context, schemaService *services.SchemaService, dao relationdbdao.IPrimaryDao) error {
	schemasDir := "erpnext_accounting/schemas"
	files, err := os.ReadDir(schemasDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Nothing to do
		}
		return fmt.Errorf("failed to read schemas directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(schemasDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read schema file %s: %w", filePath, err)
		}

		var entity descriptors.Entity
		if err := json.Unmarshal(data, &entity); err != nil {
			return fmt.Errorf("failed to parse schema file %s: %w", filePath, err)
		}

		// Check if schema already exists to avoid duplicates
		existing, err := schemaService.ByNameOrDefault(ctx, entity.Name, descriptors.EntitySchema, nil)
		if err != nil {
			return fmt.Errorf("failed to check existing schema %s: %w", entity.Name, err)
		}

		if existing != nil {
			continue // Schema already exists
		}

		// CREATE TABLE
		var cols []datamodels.Column
		// Default fields
		cols = append(cols, datamodels.Column{Name: "id", Type: datamodels.Id})
		
		for _, attr := range entity.Attributes {
			if !attr.DataType.IsLocal() {
				continue
			}
			
			var colType datamodels.ColumnType
			switch attr.DataType {
			case descriptors.Int:
				colType = datamodels.Int
			case descriptors.Float:
				colType = datamodels.Float
			case descriptors.Datetime:
				colType = datamodels.Datetime
			case descriptors.Boolean:
				colType = datamodels.Boolean
			case descriptors.Text:
				colType = datamodels.Text
			case descriptors.String:
				colType = datamodels.String
			case descriptors.DataTypeLookup:
				colType = datamodels.Int // Assuming IDs are ints here
			default:
				colType = datamodels.String
			}
			
			cols = append(cols, datamodels.Column{Name: attr.Field, Type: colType})
		}
		
		// System fields
		cols = append(cols, datamodels.Column{Name: "created_at", Type: datamodels.CreatedTime})
		cols = append(cols, datamodels.Column{Name: "updated_at", Type: datamodels.UpdatedTime})
		cols = append(cols, datamodels.Column{Name: "deleted", Type: datamodels.Boolean})
		
		err = dao.CreateTable(ctx, entity.TableName, cols)
		if err != nil && !isTableExistsError(err) {
			return fmt.Errorf("failed to create table %s: %w", entity.TableName, err)
		}

		schema := &descriptors.Schema{
			Name:              entity.Name,
			Type:              descriptors.EntitySchema,
			IsLatest:          true,
			PublicationStatus: descriptors.Published,
			Settings: &descriptors.SchemaSettings{
				Entity: &entity,
			},
		}

		_, err = schemaService.Save(ctx, schema, true)
		if err != nil {
			return fmt.Errorf("failed to save schema %s: %w", entity.Name, err)
		}
		
		fmt.Printf("Registered schema and created table: %s\n", entity.Name)
	}

	return nil
}

func isTableExistsError(err error) bool {
	// Simple check, works for sqlite usually containing "already exists"
	if err == nil {
		return false
	}
	errStr := err.Error()
	return len(errStr) > 0 && (stringContains(errStr, "already exists") || stringContains(errStr, "exists"))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}