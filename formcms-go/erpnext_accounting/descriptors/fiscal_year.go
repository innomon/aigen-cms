package descriptors

import (
	"github.com/formcms/formcms-go/core/descriptors"
	"github.com/formcms/formcms-go/utils/displaymodels"
)

var FiscalYearEntity = descriptors.Entity{
	Name:               "FiscalYear",
	DisplayName:        "Fiscal Year",
	TableName:          "fiscal_year",
	LabelAttributeName: "year",
	PrimaryKey:         "id",
	Attributes: []descriptors.Attribute{
		{
			Field:       "year",
			Header:      "Year",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Text,
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "year_start_date",
			Header:      "Start Date",
			DataType:    descriptors.Datetime,
			DisplayType: displaymodels.Date,
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "year_end_date",
			Header:      "End Date",
			DataType:    descriptors.Datetime,
			DisplayType: displaymodels.Date,
			InList:      true,
			InDetail:    true,
		},
		{
			Field:       "disabled",
			Header:      "Disabled",
			DataType:    descriptors.Boolean,
			DisplayType: displaymodels.Text,
		},
	},
}
