package descriptors

import (
	"github.com/formcms/formcms-go/core/descriptors"
	"github.com/formcms/formcms-go/utils/displaymodels"
)

var CompanyEntity = descriptors.Entity{
	Name:               "Company",
	DisplayName:        "Company",
	TableName:          "company",
	LabelAttributeName: "company_name",
	PrimaryKey:         "id",
	Attributes: []descriptors.Attribute{
		{
			Field:       "company_name",
			Header:      "Company Name",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Text,
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "abbr",
			Header:      "Abbreviation",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Text,
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "default_currency",
			Header:      "Default Currency",
			DataType:    descriptors.DataTypeLookup,
			DisplayType: displaymodels.Lookup,
			Options:     "Currency",
			InList:      true,
			InDetail:    true,
		},
	},
}
