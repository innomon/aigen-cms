package descriptors

import (
	"github.com/formcms/formcms-go/core/descriptors"
	"github.com/formcms/formcms-go/utils/displaymodels"
)

var AccountEntity = descriptors.Entity{
	Name:               "Account",
	DisplayName:        "Account",
	TableName:          "account",
	LabelAttributeName: "account_name",
	PrimaryKey:         "id",
	Attributes: []descriptors.Attribute{
		{
			Field:       "account_name",
			Header:      "Account Name",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Text,
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "account_number",
			Header:      "Account Number",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Text,
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "is_group",
			Header:      "Is Group",
			DataType:    descriptors.Boolean,
			DisplayType: displaymodels.Text, // Need to verify Boolean display type
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "company",
			Header:      "Company",
			DataType:    descriptors.DataTypeLookup,
			DisplayType: displaymodels.Lookup,
			Options:     "Company",
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "root_type",
			Header:      "Root Type",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Dropdown,
			Options:     "Asset,Liability,Income,Expense,Equity",
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "parent_account",
			Header:      "Parent Account",
			DataType:    descriptors.DataTypeLookup,
			DisplayType: displaymodels.Lookup,
			Options:     "Account",
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "account_type",
			Header:      "Account Type",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Dropdown,
			Options:     "Bank,Cash,Payable,Receivable,Stock,Tax", // Shortened list
			InList:      true,
			InDetail:    true,
		},
	},
}
