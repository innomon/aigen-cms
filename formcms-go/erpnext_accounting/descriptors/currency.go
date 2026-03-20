package descriptors

import (
	"github.com/formcms/formcms-go/core/descriptors"
	"github.com/formcms/formcms-go/utils/displaymodels"
)

var CurrencyEntity = descriptors.Entity{
	Name:               "Currency",
	DisplayName:        "Currency",
	TableName:          "currency",
	LabelAttributeName: "currency_name",
	PrimaryKey:         "id",
	Attributes: []descriptors.Attribute{
		{
			Field:       "currency_name",
			Header:      "Currency",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Text,
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "symbol",
			Header:      "Symbol",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Text,
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "fraction",
			Header:      "Fraction",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Text,
		},
		{
			Field:       "fraction_units",
			Header:      "Fraction Units",
			DataType:    descriptors.Int,
			DisplayType: displaymodels.Number,
		},
	},
}
