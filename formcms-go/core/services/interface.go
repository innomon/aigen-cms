package services

import (
	"context"

	"github.com/formcms/formcms-go/core/descriptors"
)

type ISchemaService interface {
	All(ctx context.Context, schemaType *descriptors.SchemaType, names []string, status *descriptors.PublicationStatus) ([]*descriptors.Schema, error)
	ById(ctx context.Context, id int64) (*descriptors.Schema, error)
	BySchemaId(ctx context.Context, schemaId string) (*descriptors.Schema, error)
	ByNameOrDefault(ctx context.Context, name string, schemaType descriptors.SchemaType, status *descriptors.PublicationStatus) (*descriptors.Schema, error)
	LoadEntity(ctx context.Context, name string) (*descriptors.Entity, error)
	LoadLoadedEntity(ctx context.Context, name string) (*descriptors.LoadedEntity, error)
	Save(ctx context.Context, schema *descriptors.Schema, asPublished bool) (*descriptors.Schema, error)
	Delete(ctx context.Context, schemaId string) error
}

type IEntityService interface {
	List(ctx context.Context, name string, pagination datamodels.Pagination, filters []datamodels.Filter, sorts []datamodels.Sort) ([]datamodels.Record, int64, error)
	Single(ctx context.Context, name string, id interface{}) (datamodels.Record, error)
	Insert(ctx context.Context, name string, data datamodels.Record) (datamodels.Record, error)
	Update(ctx context.Context, name string, data datamodels.Record) (datamodels.Record, error)
	Delete(ctx context.Context, name string, id interface{}) error

	CollectionList(ctx context.Context, name, id, attr string, pagination datamodels.Pagination, filters []datamodels.Filter, sorts []datamodels.Sort) ([]datamodels.Record, int64, error)
	CollectionInsert(ctx context.Context, name, id, attr string, data datamodels.Record) (datamodels.Record, error)

	JunctionList(ctx context.Context, name, id, attr string, exclude bool, pagination datamodels.Pagination, filters []datamodels.Filter, sorts []datamodels.Sort) ([]datamodels.Record, int64, error)
	JunctionSave(ctx context.Context, name, id, attr string, targetIds []interface{}) error
	JunctionDelete(ctx context.Context, name, id, attr string, targetIds []interface{}) error
}

type IGraphQLService interface {
	Query(ctx context.Context, query string, variables map[string]interface{}) (interface{}, error)
}
