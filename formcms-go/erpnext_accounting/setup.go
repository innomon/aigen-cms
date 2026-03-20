package erpnext_accounting

import (
	"context"

	"github.com/formcms/formcms-go/core/descriptors"
	"github.com/formcms/formcms-go/core/services"
	erptdesc "github.com/formcms/formcms-go/erpnext_accounting/descriptors"
)

func Setup(ctx context.Context, schemaService *services.SchemaService) error {
	entities := []descriptors.Entity{
		erptdesc.CurrencyEntity,
		erptdesc.CompanyEntity,
		erptdesc.FiscalYearEntity,
		erptdesc.CostCenterEntity,
		erptdesc.AccountEntity,
		erptdesc.GLEntryEntity,
		erptdesc.JournalEntryEntity,
		erptdesc.JournalEntryAccountEntity,
	}

	for _, entity := range entities {
		// Check if schema already exists to avoid duplicates
		existing, err := schemaService.ByNameOrDefault(ctx, entity.Name, descriptors.EntitySchema, nil)
		if err != nil {
			return err
		}
		
		if existing != nil {
			continue // Schema already exists
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
			return err
		}
	}

	return nil
}
