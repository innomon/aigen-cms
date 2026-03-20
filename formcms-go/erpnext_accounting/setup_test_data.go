package erpnext_accounting

import (
	"context"
	"fmt"

	"github.com/formcms/formcms-go/core/services"
	"github.com/formcms/formcms-go/utils/datamodels"
)

func SetupTestData(ctx context.Context, entityService services.IEntityService) error {
	// Check if data already exists
	limit := "1"
	records, _, err := entityService.List(ctx, "Currency", datamodels.Pagination{Limit: &limit}, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to list currencies: %v", err)
	}
	if len(records) > 0 {
		// Data already exists, skip
		return nil
	}

	fmt.Println("Setting up India-specific ERPNext test data...")

	// 1. Currency
	currencyRec, err := entityService.Insert(ctx, "Currency", datamodels.Record{
		"currency_name":  "INR",
		"symbol":         "₹",
		"fraction":       "Paisa",
		"fraction_units": 100,
	})
	if err != nil {
		return fmt.Errorf("failed to insert currency: %v", err)
	}
	currencyID := currencyRec["id"]

	// 2. Company
	companyRec, err := entityService.Insert(ctx, "Company", datamodels.Record{
		"company_name":     "Tata Motors",
		"abbr":             "TATA",
		"default_currency": currencyID,
	})
	if err != nil {
		return fmt.Errorf("failed to insert company: %v", err)
	}
	companyID := companyRec["id"]

	// 3. Fiscal Year
	fyRec, err := entityService.Insert(ctx, "Fiscal Year", datamodels.Record{
		"year":            "2023-2024",
		"year_start_date": "2023-04-01T00:00:00Z",
		"year_end_date":   "2024-03-31T00:00:00Z",
		"disabled":        false,
	})
	if err != nil {
		return fmt.Errorf("failed to insert fiscal year: %v", err)
	}
	// fyID := fyRec["id"]
	_ = fyRec

	// 4. Cost Centers
	ccMainRec, err := entityService.Insert(ctx, "Cost Center", datamodels.Record{
		"cost_center_name": "Main - TATA",
		"company":          companyID,
		"is_group":         true,
	})
	if err != nil {
		return fmt.Errorf("failed to insert main cost center: %v", err)
	}

	ccSalesRec, err := entityService.Insert(ctx, "Cost Center", datamodels.Record{
		"cost_center_name":   "Sales - TATA",
		"company":            companyID,
		"is_group":           false,
		"parent_cost_center": ccMainRec["id"],
	})
	if err != nil {
		return fmt.Errorf("failed to insert sales cost center: %v", err)
	}

	// 5. Accounts
	accAssetRec, err := entityService.Insert(ctx, "Account", datamodels.Record{
		"account_name": "Assets",
		"is_group":     true,
		"company":      companyID,
		"root_type":    "Asset",
	})
	if err != nil {
		return fmt.Errorf("failed to insert asset account: %v", err)
	}

	accBankGroupRec, err := entityService.Insert(ctx, "Account", datamodels.Record{
		"account_name":   "Bank Accounts",
		"is_group":       true,
		"company":        companyID,
		"root_type":      "Asset",
		"parent_account": accAssetRec["id"],
	})
	if err != nil {
		return fmt.Errorf("failed to insert bank group account: %v", err)
	}

	accHDFCRec, err := entityService.Insert(ctx, "Account", datamodels.Record{
		"account_name":   "HDFC Bank",
		"is_group":       false,
		"company":        companyID,
		"root_type":      "Asset",
		"parent_account": accBankGroupRec["id"],
		"account_type":   "Bank",
	})
	if err != nil {
		return fmt.Errorf("failed to insert HDFC account: %v", err)
	}

	accIncomeRec, err := entityService.Insert(ctx, "Account", datamodels.Record{
		"account_name": "Income",
		"is_group":     true,
		"company":      companyID,
		"root_type":    "Income",
	})
	if err != nil {
		return fmt.Errorf("failed to insert income account: %v", err)
	}

	accSalesRec, err := entityService.Insert(ctx, "Account", datamodels.Record{
		"account_name":   "Sales",
		"is_group":       false,
		"company":        companyID,
		"root_type":      "Income",
		"parent_account": accIncomeRec["id"],
		"account_type":   "Income Account",
	})
	if err != nil {
		return fmt.Errorf("failed to insert sales account: %v", err)
	}

	// 6. Journal Entry
	// Since JournalEntry has accounts via DataTypeCollection, we insert JournalEntry first
	jeRec, err := entityService.Insert(ctx, "Journal Entry", datamodels.Record{
		"voucher_type": "Journal Entry",
		"posting_date": "2023-04-05T00:00:00Z",
		"company":      companyID,
		"total_debit":  50000.0,
		"total_credit": 50000.0,
		"user_remark":  "Initial Sales Revenue",
	})
	if err != nil {
		return fmt.Errorf("failed to insert journal entry: %v", err)
	}
	jeID := jeRec["id"]

	// Insert Child Table entries via CollectionInsert
	_, err = entityService.CollectionInsert(ctx, "Journal Entry", fmt.Sprintf("%v", jeID), "accounts", datamodels.Record{
		"account":     accHDFCRec["id"],
		"debit":       50000.0,
		"credit":      0.0,
		"cost_center": ccSalesRec["id"],
		"user_remark": "Received payment",
	})
	if err != nil {
		return fmt.Errorf("failed to insert journal entry account (debit): %v", err)
	}

	_, err = entityService.CollectionInsert(ctx, "Journal Entry", fmt.Sprintf("%v", jeID), "accounts", datamodels.Record{
		"account":     accSalesRec["id"],
		"debit":       0.0,
		"credit":      50000.0,
		"cost_center": ccSalesRec["id"],
		"user_remark": "Sales booked",
	})
	if err != nil {
		return fmt.Errorf("failed to insert journal entry account (credit): %v", err)
	}

	// 7. GL Entry (Ledger impacts)
	_, err = entityService.Insert(ctx, "GL Entry", datamodels.Record{
		"posting_date": "2023-04-05T00:00:00Z",
		"account":      accHDFCRec["id"],
		"debit":        50000.0,
		"credit":       0.0,
		"voucher_type": "Journal Entry",
		"voucher_no":   fmt.Sprintf("%v", jeID),
		"company":      companyID,
		"remarks":      "Initial Sales Revenue",
	})
	if err != nil {
		return fmt.Errorf("failed to insert GL entry (debit): %v", err)
	}

	_, err = entityService.Insert(ctx, "GL Entry", datamodels.Record{
		"posting_date": "2023-04-05T00:00:00Z",
		"account":      accSalesRec["id"],
		"debit":        0.0,
		"credit":       50000.0,
		"voucher_type": "Journal Entry",
		"voucher_no":   fmt.Sprintf("%v", jeID),
		"company":      companyID,
		"remarks":      "Initial Sales Revenue",
	})
	if err != nil {
		return fmt.Errorf("failed to insert GL entry (credit): %v", err)
	}

	fmt.Println("Test data successfully created.")
	return nil
}
