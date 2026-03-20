package descriptors

import (
	"github.com/formcms/formcms-go/core/descriptors"
	"github.com/formcms/formcms-go/utils/displaymodels"
)

var CostCenterEntity = descriptors.Entity{
	Name:               "CostCenter",
	DisplayName:        "Cost Center",
	TableName:          "cost_center",
	LabelAttributeName: "cost_center_name",
	PrimaryKey:         "id",
	Attributes: []descriptors.Attribute{
		{
			Field:       "cost_center_name",
			Header:      "Cost Center Name",
			DataType:    descriptors.String,
			DisplayType: displaymodels.Text,
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
			Field:       "is_group",
			Header:      "Is Group",
			DataType:    descriptors.Boolean,
			DisplayType: displaymodels.Text,
		},
		{
			Field:       "parent_cost_center",
			Header:      "Parent Cost Center",
			DataType:    descriptors.DataTypeLookup,
			DisplayType: displaymodels.Lookup,
			Options:     "CostCenter",
		},
	},
}
