package descriptors

import (
	"github.com/formcms/formcms-go/core/descriptors"
	"github.com/formcms/formcms-go/utils/displaymodels"
)

var JournalEntryAccountEntity = descriptors.Entity{
	Name:               "JournalEntryAccount",
	DisplayName:        "Journal Entry Account",
	TableName:          "journal_entry_account",
	LabelAttributeName: "account",
	PrimaryKey:         "id",
	Attributes: []descriptors.Attribute{
		{
			Field:       "parent",
			Header:      "Journal Entry",
			DataType:    descriptors.DataTypeLookup,
			DisplayType: displaymodels.Lookup,
			Options:     "JournalEntry",
		},
		{
			Field:       "account",
			Header:      "Account",
			DataType:    descriptors.DataTypeLookup,
			DisplayType: displaymodels.Lookup,
			Options:     "Account",
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "debit",
			Header:      "Debit",
			DataType:    descriptors.Float,
			DisplayType: displaymodels.Number,
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "credit",
			Header:      "Credit",
			DataType:    descriptors.Float,
			DisplayType: displaymodels.Number,
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "cost_center",
			Header:      "Cost Center",
			DataType:    descriptors.DataTypeLookup,
			DisplayType: displaymodels.Lookup,
			Options:     "CostCenter",
		},
		{
			Field:       "user_remark",
			Header:      "User Remark",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Text,
		},
	},
}
