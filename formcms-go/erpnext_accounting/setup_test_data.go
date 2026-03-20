package erpnext_accounting

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/formcms/formcms-go/core/services"
	"github.com/formcms/formcms-go/utils/datamodels"
)

type TestDataEntry struct {
	Entity   string
	Ref      string
	Data     map[string]interface{}
	Children map[string][]map[string]interface{}
}

func SetupTestData(ctx context.Context, entityService services.IEntityService) error {
	// Check if data already exists to avoid duplication
	limit := "1"
	records, _, err := entityService.List(ctx, "Currency", datamodels.Pagination{Limit: &limit}, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to list currencies: %v", err)
	}
	if len(records) > 0 {
		return nil // Data already exists
	}

	fmt.Println("Setting up ERPNext test data from JSON...")

	dataBytes, err := os.ReadFile("erpnext_accounting/data/test_data.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Nothing to do
		}
		return fmt.Errorf("failed to read test_data.json: %v", err)
	}

	var entries []TestDataEntry
	if err := json.Unmarshal(dataBytes, &entries); err != nil {
		return fmt.Errorf("failed to unmarshal test data: %v", err)
	}

	refMap := make(map[string]interface{})

	resolveRefs := func(data map[string]interface{}) {
		for k, v := range data {
			if strVal, ok := v.(string); ok && strings.HasPrefix(strVal, "$Ref:") {
				refKey := strings.TrimPrefix(strVal, "$Ref:")
				if resolvedVal, exists := refMap[refKey]; exists {
					data[k] = resolvedVal
				} else {
					fmt.Printf("Warning: Could not resolve reference %s\n", refKey)
				}
			}
		}
	}

	for _, entry := range entries {
		// Prepare data, resolve references
		resolveRefs(entry.Data)

		rec, err := entityService.Insert(ctx, entry.Entity, entry.Data)
		if err != nil {
			return fmt.Errorf("failed to insert %s (Ref: %s): %v", entry.Entity, entry.Ref, err)
		}

		if entry.Ref != "" {
			refMap[entry.Ref] = rec["id"]
		}

		// Insert children
		if entry.Children != nil {
			for childAttr, childrenArr := range entry.Children {
				for i, childData := range childrenArr {
					resolveRefs(childData)
					_, err = entityService.CollectionInsert(ctx, entry.Entity, fmt.Sprintf("%v", rec["id"]), childAttr, childData)
					if err != nil {
						return fmt.Errorf("failed to insert child %s for %s (Ref: %s) at index %d: %v", childAttr, entry.Entity, entry.Ref, i, err)
					}
				}
			}
		}
	}

	fmt.Println("Test data successfully created from JSON.")
	return nil
}