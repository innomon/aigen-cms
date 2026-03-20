package descriptors

import (
	"github.com/formcms/formcms-go/core/descriptors"
	"github.com/formcms/formcms-go/utils/displaymodels"
)

var JournalEntryEntity = descriptors.Entity{
	Name:               "JournalEntry",
	DisplayName:        "Journal Entry",
	TableName:          "journal_entry",
	LabelAttributeName: "voucher_type",
	PrimaryKey:         "id",
	Attributes: []descriptors.Attribute{
		{
			Field:       "voucher_type",
			Header:      "Entry Type",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Dropdown,
			Options:     "Journal Entry,Inter Company Journal Entry,Bank Entry,Cash Entry,Credit Card Entry,Debit Note,Credit Note,Contra Entry,Excise Entry,Write Off Entry,Opening Entry,Depreciation Entry,Asset Disposal,Periodic Accounting Entry,Exchange Rate Revaluation,Exchange Gain Or Loss,Deferred Revenue,Deferred Expense",
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "posting_date",
			Header:      "Posting Date",
			DataType:    descriptors.Datetime,
			DisplayType: displaymodels.Date,
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
			Field:       "accounts",
			Header:      "Accounting Entries",
			DataType:    descriptors.DataTypeCollection,
			DisplayType: displaymodels.EditTable,
			Options:     "JournalEntryAccount|parent",
			InDetail:    true,
		},
		{
			Field:       "total_debit",
			Header:      "Total Debit",
			DataType:    descriptors.Float,
			DisplayType: displaymodels.Number,
			InList:      true,
		},
		{
			Field:       "total_credit",
			Header:      "Total Credit",
			DataType:    descriptors.Float,
			DisplayType: displaymodels.Number,
		},
		{
			Field:       "user_remark",
			Header:      "User Remark",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Textarea,
		},
		{
			Field:       "cheque_no",
			Header:      "Reference Number",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Text,
		},
		{
			Field:       "cheque_date",
			Header:      "Reference Date",
			DataType:    descriptors.Datetime,
			DisplayType: displaymodels.Date,
		},
	},
}
