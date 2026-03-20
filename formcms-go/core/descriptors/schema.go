package descriptors

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

type SchemaType string

const (
	MenuSchema   SchemaType = "Menu"
	EntitySchema SchemaType = "Entity"
	QuerySchema  SchemaType = "Query"
	PageSchema   SchemaType = "Page"
)

type SchemaSettings struct {
	Entity *Entity `json:"entity,omitempty" mapstructure:"entity"`
	Query  *Query  `json:"query,omitempty" mapstructure:"query"`
	Menu   *Menu   `json:"menu,omitempty" mapstructure:"menu"`
	Page   *Page   `json:"page,omitempty" mapstructure:"page"`
}

type Schema struct {
	Id                int64             `json:"id" mapstructure:"id"`
	SchemaId          string            `json:"schemaId" mapstructure:"schemaId"`
	Name              string            `json:"name" mapstructure:"name"`
	Type              SchemaType        `json:"type" mapstructure:"type"`
	Settings          *SchemaSettings   `json:"settings" mapstructure:"settings"`
	Description       string            `json:"description" mapstructure:"description"`
	IsLatest          bool              `json:"isLatest" mapstructure:"isLatest"`
	PublicationStatus PublicationStatus `json:"publicationStatus" mapstructure:"publicationStatus"`
	CreatedAt         time.Time         `json:"createdAt" mapstructure:"createdAt"`
	CreatedBy         string            `json:"createdBy" mapstructure:"createdBy"`
	Deleted           bool              `json:"deleted" mapstructure:"deleted"`
}

func RecordToSchema(record map[string]interface{}) (*Schema, error) {
	var s Schema
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &s,
		TagName:  "mapstructure",
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			func(f interface{}, t interface{}) (interface{}, error) {
				if f == nil {
					return nil, nil
				}
				// Handle string to json.RawMessage or struct conversion for Settings
				if str, ok := f.(string); ok {
					// Check if target is *SchemaSettings
					if fmt.Sprintf("%T", t) == "*descriptors.SchemaSettings" {
						var settings SchemaSettings
						if err := json.Unmarshal([]byte(str), &settings); err != nil {
							return nil, err
						}
						return &settings, nil
					}
				}
				return f, nil
			},
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(record); err != nil {
		return nil, err
	}

	return &s, nil
}
