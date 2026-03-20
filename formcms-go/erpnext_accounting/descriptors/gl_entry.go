package descriptors

import (
	"github.com/formcms/formcms-go/core/descriptors"
	"github.com/formcms/formcms-go/utils/displaymodels"
)

var GLEntryEntity = descriptors.Entity{
	Name:               "GLEntry",
	DisplayName:        "GL Entry",
	TableName:          "gl_entry",
	LabelAttributeName: "voucher_no",
	PrimaryKey:         "id",
	Attributes: []descriptors.Attribute{
		{
			Field:       "posting_date",
			Header:      "Posting Date",
			DataType:    descriptors.Datetime,
			DisplayType: displaymodels.Date,
			InList:      true,
			InDetail:    true,
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
			Field:       "voucher_type",
			Header:      "Voucher Type",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Text,
			InList:      true,
		},
		{
			Field:       "voucher_no",
			Header:      "Voucher No",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Text,
			InList:      true,
		},
		{
			Field:       "company",
			Header:      "Company",
			DataType:    descriptors.DataTypeLookup,
			DisplayType: displaymodels.Lookup,
			Options:     "Company",
			InList:      true,
		},
		{
			Field:       "remarks",
			Header:      "Remarks",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Textarea,
		},
	},
}
